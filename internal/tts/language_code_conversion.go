package tts

// List of supported languages
var LanguageMap = map[string]string{
	"eng":    "en", // English US
	"eng-UK": "en-UK",
	"eng-AU": "en-AU",
	"jpn":    "ja", // Japanese
	"deu":    "de", // German
	"spa":    "es", // Spanish
	"rus":    "ru", // Russian
	"ara":    "ar", // Arabic
	"ben":    "bn", // Bengali
	"ces":    "cs", // Czech
	"dan":    "da", // Danish
	"nld":    "nl", // Dutch
	"fin":    "fi", // Finnish
	"ell":    "el", // Greek
	"hin":    "hi", // Hindi
	"hun":    "hu", // Hungarian
	"ind":    "id", // Indonesian
	"khm":    "km", // Khmer
	"lat":    "la", // Latin
	"ita":    "it", // Italian
	"nor":    "no", // Norwegian (nb)
	"pol":    "pl", // Polish
	"slk":    "sk", // Slovak
	"swe":    "sv", // Swedish
	"tha":    "th", // Thai
	"tur":    "tr", // Turkish
	"ukr":    "uk", // Ukrainian
	"vie":    "vi", // Vietnamese
	"afr":    "af", // Afrikaans
	"bul":    "bg", // Bulgarian
	"cat":    "ca", // Catalan
	"cym":    "cy", // Welsh
	"est":    "et", // Estonian
	"fra":    "fr", // French
	"guj":    "gu", // Gujarati
	"isl":    "is", // Icelandic
	"jav":    "jv", // Javanese
	"kan":    "kn", // Kannada
	"kor":    "ko", // Korean
	"lav":    "lv", // Latvian
	"mal":    "ml", // Malayalam
	"mar":    "mr", // Marathi
	"msa":    "ms", // Malay
	"nep":    "ne", // Nepali
	"por":    "pt", // Portuguese
	"ron":    "ro", // Romanian
	"sin":    "si", // Sinhala
	"srp":    "sr", // Serbian
	"sun":    "su", // Sundanese
	"tam":    "ta", // Tamil
	"tel":    "te", // Telugu
	"tgl":    "tl", // Tagalog
	"urd":    "ur", // Urdu
	"zho":    "zh", // Chinese
	"swa":    "sw", // Swahili
	"sqi":    "sq", // Albanian
	"mya":    "my", // Burmese
	"mkd":    "mk", // Macedonian
	"hye":    "hy", // Armenian
	"hrv":    "hr", // Croatian
	"epo":    "eo", // Esperanto
	"bos":    "bs", // Bosnian
}

// This function coverts the ISO 639-3 language code used by the API.Bible API to the format used by the text-to-speech module.
func ISOtoTTScode(ISO string) string {
	result, ok := LanguageMap[ISO]
	if !ok {
		return ""
	}
	return result
}
