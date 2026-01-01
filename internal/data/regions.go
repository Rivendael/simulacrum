package data

// AddressRegion groups cities, states, and countries for deterministic generation
type AddressRegion struct {
	Country string
	States  []string
	Cities  []string
}

// AddressRegions organized by geographic region to ensure consistency
var AddressRegions = []AddressRegion{
	{
		Country: "USA",
		States:  []string{"AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DE", "FL", "GA", "HI", "ID", "IL", "IN", "IA", "KS", "KY", "LA", "ME", "MD", "MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ", "NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA", "RI", "SC", "SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY"},
		Cities:  []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose", "Austin", "Jacksonville", "Fort Worth", "Columbus", "Indianapolis", "Charlotte", "San Francisco", "Seattle", "Denver", "Boston", "Memphis", "Nashville", "Detroit", "Oklahoma City", "Portland", "Las Vegas", "Louisville", "Baltimore", "Milwaukee", "Albuquerque", "Tucson", "Fresno", "Long Beach", "Kansas City", "Mesa", "Atlanta", "Miami", "Arlington", "New Orleans", "Bakersfield", "Tampa", "Aurora", "Anaheim", "Santa Ana", "Riverside", "Corpus Christi", "Lexington", "Henderson", "Plano", "Stockton", "St. Louis"},
	},
	{
		Country: "Canada",
		States:  []string{"AB", "BC", "MB", "NB", "NL", "NS", "ON", "PE", "QC", "SK"},
		Cities:  []string{"Toronto", "Montreal", "Vancouver", "Calgary", "Edmonton", "Ottawa", "Winnipeg", "Quebec City", "Hamilton", "Kitchener", "London", "Halifax", "Windsor", "Saskatoon", "Laval", "Victoria", "Barrie", "St. Catharines", "Markham", "Mississauga"},
	},
	{
		Country: "Mexico",
		States:  []string{"AGS", "BC", "BCS", "CAM", "COAH", "COL", "CDMX", "DGO", "GTO", "GRO", "HGO", "JAL", "MEX", "MICH", "MOR", "NAY", "OAX", "PUE", "QRO", "QROO", "SLP", "SIN", "SON", "TAB", "TAMPS", "TLAX", "VER", "YUC", "ZAC"},
		Cities:  []string{"Mexico City", "Guadalajara", "Monterrey", "Ecatepec", "Puebla", "Toluca", "Leon", "Cancun", "Irapuato", "Juarez", "Zapopan", "Chihuahua", "Morelia", "Hermosillo", "Saltillo", "Merida", "Aguascalientes", "Veracruz", "Culiacan", "Celaya"},
	},
	{
		Country: "France",
		States:  []string{"75", "78", "91", "92", "93", "94", "95", "13", "69", "67", "59"},
		Cities:  []string{"Paris", "Marseille", "Lyon", "Toulouse", "Nice", "Nantes", "Strasbourg", "Montpellier", "Bordeaux", "Lille", "Rennes"},
	},
	{
		Country: "Germany",
		States:  []string{"BW", "BY", "BE", "BB", "HB", "HH", "HE", "MV", "NI", "NW", "RP", "SL", "SN", "ST", "SH", "TH"},
		Cities:  []string{"Berlin", "Munich", "Cologne", "Frankfurt", "Hamburg", "Dusseldorf", "Stuttgart", "Dortmund", "Essen", "Leipzig", "Dresden", "Hanover"},
	},
	{
		Country: "Japan",
		States:  []string{"TO", "KA", "OS", "HY", "SA", "SH", "AIT", "GIF", "SZ", "CHB", "TOK"},
		Cities:  []string{"Tokyo", "Yokohama", "Osaka", "Kobe", "Kyoto", "Kawasaki", "Saitama", "Hiroshima", "Fukuoka", "Nagoya", "Sapporo"},
	},
	{
		Country: "United Kingdom",
		States:  []string{"ENG", "SCT", "WAL", "NIR"},
		Cities:  []string{"London", "Manchester", "Birmingham", "Leeds", "Glasgow", "Sheffield", "Bristol", "Edinburgh", "Liverpool", "York", "Cambridge"},
	},
	{
		Country: "Spain",
		States:  []string{"MA", "BA", "CA", "CM", "CL", "CT", "VC", "GA", "LR", "NA", "AR", "AS", "CN", "CB", "CE", "EX", "MD", "ME", "MU", "PM"},
		Cities:  []string{"Madrid", "Barcelona", "Valencia", "Seville", "Bilbao", "Malaga", "Murcia", "Palma", "Las Palmas", "Alicante", "Cordoba"},
	},
	{
		Country: "Italy",
		States:  []string{"AG", "AL", "AN", "AO", "AR", "AT", "AV", "BA", "BL", "BN", "BR", "BS", "BZ", "CA", "CB", "CE", "CH", "CL", "CN", "CO", "CS", "CT", "CZ", "EN", "EX", "FC", "FE", "FG", "FI", "FM", "FR", "GE", "GO", "GR", "IM", "IS", "KR", "LC", "LE", "LI", "LO", "LT", "LU", "MB", "MC", "ME", "MI", "MN", "MO", "MS", "MT", "NA", "NO", "NU", "OG", "OR", "OT", "PA", "PC", "PD", "PE", "PG", "PI", "PN", "PO", "PR", "PS", "PT", "PU", "PZ", "RA", "RC", "RE", "RG", "RI", "RM", "RN", "RO", "SA", "SI", "SO", "SP", "SS", "SV", "TA", "TE", "TN", "TO", "TP", "TR", "TS", "TT", "TV", "UA", "UD", "VA", "VB", "VC", "VE", "VI", "VR", "VS"},
		Cities:  []string{"Rome", "Milan", "Naples", "Turin", "Palermo", "Genoa", "Bologna", "Florence", "Bari", "Catania", "Venice"},
	},
}
