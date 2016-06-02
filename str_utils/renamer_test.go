package strutils

import (
	"regexp"
	"testing"
)

type TestFixture struct {
	typicalNames       []string
	randomSimilarStrs1 []string
	randomSimilarStrs2 []string
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

	re := regexp.MustCompile("(S?\\d{1,2})(E?\\d{2})")
	stripedStrs := RemoveCommonSubstringsPreseveMatch(episodes, 0.5, re)
	if stripedStrs[0].New_str != "Anthropology.101" {
		t.Errorf("%s didn't match %s", stripedStrs[0].New_str, "Anthropology.101")
	}
}
