package types

var (
	B1     = Base{}
	C1s2eA = Subsub{
		Title:   "Chap1Sub2ExtA",
		Path:    "extA",
		Source:  "../../misc/test/sub4",
		Weight:  10,
		Flavour: "eng",
		Enum:    "a. ",
	}
	C1s1 = Subchapter{
		Title:   "Chap1Sub1",
		Path:    "sub1",
		Source:  "../../misc/test/sub1",
		Weight:  10,
		Flavour: "eng",
		Enum:    "1. ",
		Subsubs: []Subsub{},
	}
	C1s2 = Subchapter{
		Title:   "Chap1Sub2",
		Path:    "sub2",
		Source:  "../../misc/test/sub2",
		Weight:  20,
		Flavour: "eng",
		Enum:    "2. ",
		Subsubs: []Subsub{C1s2eA},
	}
	C1 = Chapter{
		Title:    "Chap1",
		Path:     "chap1",
		Source:   "../../misc/test/chap1",
		Weight:   10,
		Flavour:  "eng",
		Enum:     "I. ",
		Subchaps: []Subchapter{C1s1},
	}
	C12 = Chapter{
		Title:    "Chap1",
		Path:     "chap1",
		Source:   "../../misc/test/chap1",
		Weight:   10,
		Flavour:  "eng",
		Enum:     "I. ",
		Subchaps: []Subchapter{C1s1, C1s2},
	}
)
