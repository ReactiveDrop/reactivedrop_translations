#include <iostream>
#include <fstream>
#include <string>
#include <regex>
#include <vector>

using namespace std;

vector<string> languageList = {
	"brazilian",
	"czech",
	"danish",
	"dutch",
	"english",
	"finnish",
	"french",
	"german",
	"hungarian",
	"indonesian",
	"italian",
	"japanese",
	"koreana",
	"latam",
	"norwegian",
	"polish",
	"portuguese",
	"romanian",
	"russian",
	"schinese",
	"spanish",
	"swedish",
	"tchinese",
	"thai",
	"turkish",
	"ukrainian",
	"vietnamese" 
};


int main()
{
	string filePath = "C:\\Users\\ustcy\\Desktop\\reactivedrop_translations\\addons\\traitors_challenge\\resource\\";
	ofstream outFileStream(filePath + "AddonTranslationTool\\translations_all.nut");
	if (!outFileStream.is_open())
	{
		cout << "Error: " << filePath + "reactivedrop_all.txt" << " not found." << endl;
		return 0;
	}
	outFileStream << "//This is a file that will be merged into add-on's vscript.\n//DO NOT translate this file.\n\n";

	string translatedLang;
	string translatedName;
	for (string language : languageList)
	{
		string path = filePath + "reactivedrop_" + language + ".txt";
		ifstream inFileSteam(path);
		if (!inFileSteam.is_open())
		{
			cout << "Error: " << path << " not found." << endl;
			continue;
		}
		string line;
		while (getline(inFileSteam, line))
		{
			regex regName(".*\"challenge_traitors_name\".*?\"(.*)\"");
			smatch matchName;
			if (regex_match(line, matchName, regName))
			{
				string str = matchName[1].str();
				if (!str._Equal("Traitors!"))
				{
					translatedLang = translatedLang + language + " / ";
					translatedName = translatedName + str + " / ";
				}
			}
		}
		inFileSteam.close();
	}
	outFileStream << "/*------------------------\nTranslated: " << translatedLang << endl << translatedName << "\n------------------------*/" << endl;

	outFileStream << "g_localizations <- {\n";
	for (string language : languageList)
	{
		string filePathName = filePath + "reactivedrop_" + language + ".txt";
		ifstream inFileStream(filePathName);
		if (!inFileStream.is_open())
		{
			cout << "Error: " << filePathName << " not found." << endl;
			continue;
		}

		outFileStream << "\t" << language << " = {\n";
		string line;
		while (getline(inFileStream, line))
		{
			regex reg(".*\"(challenge_.*?)\".*?\"(.*)\"");
			smatch match;
			if (regex_match(line, match, reg))
			{
				outFileStream << "\t\t" << match[1] << " = \"" << match[2] << "\",\n";
			}
		}
		outFileStream << "\t},\n";
		inFileStream.close();
	}
	outFileStream << "};";
	outFileStream.close();
	return 0;
}