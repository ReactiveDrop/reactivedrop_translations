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

function getPlaceHolders(&$v): array {
	// XXX: test string
	// $v = "<50 % %1s %+slot1% \n\n \t <B><I><clr:255:1:0>";

	$placeholders = [];

	// try isolate single tokens
	$matches = [];
	preg_match_all("/\\n|\\t|\\r|<[\w:\/]+>|%\d\w|\.[A-Z]{3}/si", $v, $matches);

	// add placeholders
	$placeholders = $matches[0];

	// isolate range patterns
	preg_match_all("/%[\w+-_]+%/s", $v, $matches);
	foreach ($matches[0] as $m) $placeholders[] = $m;

	// replace them
	foreach ($placeholders as $n => $p) {
		$v = str_replace($p, sprintf('%%#%d%%', $n), $v);
	}

	return $placeholders;
}

$languages = json_decode(
	file_get_contents(__DIR__ . '/language.mapping.json')
);

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
					if (!strstr($k, '[english]') && $originals[$k] === $v && trim($v)) {

						if ($progress->{$item->getBasename()} <= $i) {

							$placeholders = getPlaceHolders($v);

							// try translation
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
								echo sprintf('[%s] %s', $e->getMessage(), $v) . PHP_EOL;
							}

							// restore placeholders
							foreach ($placeholders as $n => $p) {
								$find = sprintf('/%% ?#%d%%/s', $p);
								$v = preg_replace($find, $p, $v);
							}
						}

						if ($i > $progress->{$item->getBasename()}) {
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
