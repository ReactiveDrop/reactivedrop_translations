#!/usr/bin/php
<?php
use Jefs42\LibreTranslate;
use Mintopia\VDFKeyValue\Encoder;
use VdfParser\Parser;
use Dariuszp\CliProgressBar;

require(__DIR__ . '/vendor/autoload.php');

// progress file
$progressFile = './.translation.progress.json';

// open progress file
$progress = file_exists($progressFile) ? json_decode(file_get_contents($progressFile)) : new stdClass;

// argc
if (in_array('--force', $argv)) $progress = new StdClass;

// languages
$languages = json_decode(
	file_get_contents(__DIR__ . '/language.mapping.json')
);

if (in_array('--lang', $argv)) {
	$pos = array_search('--lang', $argv);
	$lang = $argv[$pos+1];
	$languages = (object)[
		$lang => $languages->{$lang}
	];
}

// debug
function dd($msg)
{
	var_dump($msg);
	die;
}

function writeProgress()
{
	global $progressFile, $progress;

	file_put_contents($progressFile, json_encode($progress));
}

// mark special chars
$chars = [
	"\n" => "<br/>",
	"\t" => "<blockquote/>"
];

function escapeString($v): string {
	global $chars;

	// placeholders, .extensions and hostnames
	$v = preg_replace('/(%\d\w|\.[A-Z]{3}|< \d+|[a-z0-9:.]{12,}|%[\w+-_]+%|\s%\s)/', '<code x-data="$1"></code>', $v);

	// chars
	$v = str_replace(array_keys($chars), array_values($chars), $v);

	return $v;
}


// libs
$parser = new Parser;
$encoder = new Encoder;

$translator = new LibreTranslate(
	getenv('LIBRETRANSLATE_HOST') ?? 'http://localhost',
	getenv('LIBRETRANSLATE_PORT') ?? 5000
);

$key = getenv('LIBRETRANSLATE_API_KEY');
if ($key) {
	$translator->setApiKey($key);
}

// traverse the repository
$dir = new RecursiveDirectoryIterator(__DIR__ . '/..');
$iterator = new RecursiveIteratorIterator($dir);

foreach ($iterator as $item) {
	// only translate kv files
	if ($item->isFile() && $item->getExtension() === 'txt') {
		$p = preg_split('/[_\.]/', $item->getBasename());
		end($p);
		$lang = prev($p);

		// only auto translate non-english text
		$l = isset($languages->{$lang}) ? strlen($languages->{$lang}) : 0;
		if ($l === 2 || $l === 5) {

			// abort
			$abort = false;

			// open file
			$data = $parser->parse(
				file_get_contents($item->getPathname())
			);

			// tokens
			$originals = [];

			// progress
			if (!isset($progress->{$item->getBasename()})) $progress->{$item->getBasename()} = 0;

			// iterate data, find original text first
			if (isset($data['lang']['Tokens'])) {

				$tokens =& $data['lang']['Tokens'];
				$bar = new CliProgressBar(count($tokens), 0, sprintf("%32s", $item->getBasename()));
				$bar->display();

				foreach ($tokens as $k => $v) {
					if (strstr($k, '[english]')) {
						$key = str_replace('[english]', '', $k);
						$originals[$key] = $v;
					}
				}

				// iterate again, to find untranslated text
				$i = 0;
				foreach ($tokens as $k => &$v) {

					$i++;

					// if translation is set, and is the same, we probably need an update
					if (!$abort && !strstr($k, '[english]') && $originals[$k] === $v && trim($v)) {

						if ($progress->{$item->getBasename()} <= $i) {

							// XXX: test strings
							$test = null;
							$test = "< 50 % %1s %+slot1% \n\n \t <B><I><clr:255:1:0>";
							// $test = "< 50";
							// $test = "Examples:\ntfc.valvesoftware.com\ncounterstrike.speakeasy.net:27016\n205.158.143.200:27015";

							// try translation
							if ($test) $v = $test;

							// placeholders
							$v = escapeString($v);
							$vi = $v;

							try {
								$translation = $translator->translate($v, "auto", $languages->{$lang});

								if ($translation && $translation !== $v) {

									// if the original translation wasn't ending on a dot ., revert this
									$dot = '/\.$/';
									if (preg_match($dot, $translation) && !preg_match($dot, $v)) {
										$translation = preg_replace($dot, '', $translation);
									}

									// write back to kv, but mark translation with a *
									$v = sprintf('%s*', $translation);
								}

							} catch (Exception $e) {
								echo PHP_EOL;
								echo sprintf('error : %s', $e->getMessage()) . PHP_EOL;
								echo sprintf('text  : %s', $v);

								if (strstr($e->getMessage(), 'is not supported')) $abort = true;
							}

							// if the response was an empty string, restore original string
							$vt = $v;

							if ($v === '_*') {
								$v = $vi;
							} else {
								// strip all tags and put back chars
								$v = str_replace(
									array_values($chars),
									array_keys($chars),
									$v
								);

								// restore placeholders
								$v = preg_replace('/<code x-data="([^"]*?)"><\/code>/s', '$1', $v);

								// restore specific html tags
								$tags = ['b', 'i', 'u'];
								foreach ($tags as $tag) {
									$v = str_replace(sprintf('<%s>', $tag), sprintf('<%s>', strtoupper($tag)), $v);
									$v = str_replace(sprintf('</%s>', $tag), sprintf('</%s>', strtoupper($tag)), $v);
								}
							}

							if ($test) {
								echo PHP_EOL;
								echo sprintf('original     : %s', $test) . PHP_EOL;
								echo sprintf('input        : %s', $vi) . PHP_EOL;
								echo sprintf('output       : %s', $vt) . PHP_EOL;
								echo sprintf('translated   : %s', $v) . PHP_EOL;
								exit;
							}
						}

						if (!$abort && $i > $progress->{$item->getBasename()}) {
							$progress->{$item->getBasename()} = $i;
							if (++$i % 1000 == 0) {
								writeProgress();
							}
						}
					}

					$bar->progress();
				}

				// encode back to file
				$encoded = $encoder->encode('lang', $data['lang']);

				if ($encoded) {
					// strip indentation, RD doesn't use indents
					$encoded = preg_replace("/\"\t+\"/", "\"\t\t\"", $encoded);

					file_put_contents($item->getPathname(), $encoded);
				}

				$bar->end();
				writeProgress();

				if (in_array('--test', $argv)) die('eof');
			}
		}
	}
}


echo "+ all finished";
