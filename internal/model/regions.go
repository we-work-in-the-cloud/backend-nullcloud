package model

var ValidRegions = []string{
	"us-east", "us-west", "us-central", "eu-central", "eu-east", "eu-west",
}

type Zone struct {
	Name string `json:"name"`
}

type Region struct {
	Name  string `json:"name"`
	Zones []Zone `json:"zones"`
}

func Regions() []Region {
	regions := make([]Region, len(ValidRegions))
	for i, r := range ValidRegions {
		regions[i] = Region{
			Name: r,
			Zones: []Zone{
				{Name: r + "-1"},
				{Name: r + "-2"},
				{Name: r + "-3"},
			},
		}
	}
	return regions
}

func IsValidRegion(region string) bool {
	for _, r := range ValidRegions {
		if r == region {
			return true
		}
	}
	return false
}

// IsValidZoneForRegion reports whether zone is one of the three zones within region.
func IsValidZoneForRegion(region, zone string) bool {
	return zone == region+"-1" || zone == region+"-2" || zone == region+"-3"
}
