#!/usr/bin/php
<?php
use Jefs42\LibreTranslate;
use Mintopia\VDFKeyValue\Encoder;
use VdfParser\Parser;
use Dariuszp\CliProgressBar;

require(__DIR__ . '/vendor/autoload.php');

// debug
function dd($msg) {
	var_dump($msg);
	die;
}

// languages
$languages = [
	'french' => 'fr'
];

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
		if (isset($languages[$lang])) {

			// open file
			$data = $parser->parse(
				file_get_contents($item->getPathname())
			);

			// tokens
			$originals = [];

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
				foreach ($tokens as $k => &$v) {
					// if translation is set, and is the same, we probably need an update
					if (!strstr($k, '[english]') && $originals[$k] === $v && trim($v)) {

						// try translation
						try {
							$translation = $translator->translate($v, "auto", $languages[$lang]);

							if ($translation && $translation !== $v) {
								// write back to kv, but mark translation with a *
								$v = $translation . '*';
							}
						} catch (Exception $e) {
							echo sprintf('[%s] %s', $e->getMessage(), $v) . PHP_EOL;
						}
					}

					$bar->progress();
				}

				// encode back to file
				$encoded = $encoder->encode('lang', $data['lang']);

				if ($encoded) {
					// strip indentation, RD doesn't use indents
					$encoded = preg_replace("/^\s+/m", '', $encoded);

					file_put_contents($item->getPathname(), $encoded);
				}

				$bar->end();
			}
		}
	}
}

echo "+ all finished";
