package strutils

import (
	// "path/filepath"
	"regexp"
	"strings"
	"testing"
)

type TestFixture struct {
	typicalNames          []string
	typicalNamesPreserved []string
	randomSimilarStrs1    []string
	randomSimilarStrs2    []string
}

func getSharedStructures() TestFixture {
	fixture := TestFixture{
		typicalNames: []string{
			"Community.S02E01.Anthropology.101.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E02.Accounting.For.Lawyers.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E03.The.Psychology.Of.Letting.Go.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E04.Basic.Rocket.Science.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E05.Messianic.Myths.And.Ancient.Peoples.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E06.Epidemiology.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E07.Aerodynamics.Of.Gender.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E08.Cooperative.Calligraphy.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E09.Conspiracy.Theories.And.Interior.Design.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E10.Mixology.Certification.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E11.Abeds.Uncontrollable.Christmas.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E12.Asian.Population.Studies.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E13.Celebrity.Pharmacology.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E14.Advanced.Dungeons.And.Dragons.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E15.Early.21st.Century.Romanticism.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E16.Intermediate.Documentary.Filmmaking.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E17.Intro.To.Political.Science.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E18.Custody.Law.And.Eastern.European.Diplomacy.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E19.Critical.Film.Studies.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E20.Competitive.Wine.Tasting.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E21.Paradigms.Of.Human.Memory.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E22.Applied.Anthropology.And.Culinary.Arts.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E23.A.Fistful.Of.Paintballs.Part.1.720p.WEB-DL.x264-LeRalouf.mkv",
			"Community.S02E24.For.A.Few.Paintballs.More.Part.2.720p.WEB-DL.x264-LeRalouf.mkv",
		},
		typicalNamesPreserved: []string{
			"S02E01 Anthropology 101.mkv",
			"S02E02 Accounting For Lawyers.mkv",
			"S02E03 The Psychology Of Letting Go.mkv",
			"S02E04 Basic Rocketr Science.mkv",
			"S02E05 Messianic Myths And Ancient Peoples.mkv",
			"S02E06 Epidemiology.mkv",
			"S02E07 Aerodynamics Of Gender.mkv",
			"S02E08 Cooperative Calligraphy mkv",
			"S02E09 Conspiracy Theories And Interior Design.mkv",
			"S02E10 Mixology Certification.mkv",
			"S02E11 Abeds Uncontrollable Christmas.mkv",
			"S02E12 Asian Population Studies.mkv",
			"S02E13 Celebrity Pharmacology.mkv",
			"S02E14 Advanced Dungeons And Dragons.mkv",
			"S02E15 Early 21st Century Romanticism.mkv",
			"S02E16 Intermediate Documentary Filmmaking.mkv",
			"S02E17 Intro To Political Science.mkv",
			"S02E18 Custody Law And Eastern European Diplomacy.mkv",
			"S02E19 Critical Film Studies.mkv",
			"S02E20 Competitive Wine Tasting.mkv",
			"S02E21 Paradigms Of Human Memory.mkv",
			"S02E22 Applied Anthropology And Culinary Arts.mkv",
			"S02E23 A Fistful Of Paintballs Part 1.mkv",
			"S02E24 For A Few Paintballs More Part 2.mkv",
		},
		randomSimilarStrs1: []string{
			"catmanplan",
			"hatmanplan",
			"shatmanplan",
			"fatmanplan",
		},
		randomSimilarStrs2: []string{
			"Hello my name is Erik and I am a fan of cheese",
			"Hello my name is bill and I am a fan of popcorn",
			"Hello my name is bobby and I am a fan of cats",
			"I like food",
		},
	}
	return fixture
}

func TestFindCommonSubStrs(t *testing.T) {
	var strs = getSharedStructures().randomSimilarStrs1

	stripedStrs := RemoveCommonSubstrings(strs, 0.5)
	if stripedStrs[0].New_str != "c" {
		t.Errorf("%s didn't match %s", stripedStrs[0].New_str, "c")
	}
}

func TestFindCommonSubStrsPreserveMatch(t *testing.T) {
	var episodes = getSharedStructures().typicalNames
	var expectedOutput = getSharedStructures().typicalNamesPreserved

	// A better regex would include the file extention like so (S\d{1,2})(E\d{2})(.*)(\.mkv)
	// but what happens with the third match? The algorithm on the backend would
	// have to be changed
	re := regexp.MustCompile("(S?\\d{1,2})(E?\\d{2})")
	stripedStrs := RemoveCommonSubstringsPreseveMatch(episodes, 0.5, re)

	for index, correct_str := range expectedOutput {
		// index is the index where we are
		// extension := filepath.Ext(stripedStrs[index].New_str)
		// name := strings.TrimSuffix(stripedStrs[index].New_str, extension))
		clean := strings.Replace(stripedStrs[index].New_str, ".", " ", strings.Count(stripedStrs[index].New_str, ".")-1)

		if correct_str != stripedStrs[index].New_str {
			t.Errorf("%s didn't match %s", clean, correct_str)
		}
	}

	if stripedStrs[0].New_str != "Anthropology.101" {
		t.Errorf("%s didn't match %s", stripedStrs[0].New_str, "Anthropology.101")
	}
}
