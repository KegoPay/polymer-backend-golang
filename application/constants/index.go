package constants

const (
	SUPPORT_EMAIL string = "support@usepolymer.co"
	MAXIMUM_INTERNATIONAL_TRANSFER_LIMIT = 20000.0 // 20k dollars
	MAXIMUM_LOCAL_TRANSFER_LIMIT = 10000000 // 10M naira
	MINIMUM_INTERNATIONAL_TRANSFER_LIMIT = 1.0 // 1 dollar
	MINIMUM_INTERNATIONAL_MOMO_TRANSFER_LIMIT = 3.0 // 3 dollars
	MINIMUM_LOCAL_TRANSFER_LIMIT = 10000 // 100 naira
	BUSINESS_WALLET_LIMIT int = 1
	MAX_TRANSACTION_PIN_TRIES  int = 5
	MAX_PASSWORD_TRIES  int = 5
	INTERNATIONAL_TRANSACTION_FEE_RATE float32 = 0.01
	INTERNATIONAL_PROCESSOR_FEE_RATE float32 = 0.005
	LOCAL_TRANSACTION_FEE_VAT float32 = 0.075
	LOCAL_TRANSACTION_FEE_RATE float32 = 0.5
	LOCAL_PROCESSOR_FEE_LT_5000 float32 = 10.00
	LOCAL_PROCESSOR_FEE_LT_50000 float32 = 25.00
	LOCAL_PROCESSOR_FEE_GT_50000 float32 = 50.00

	// any new polymer account is on tier1
	TIER_ONE_DAILY_TRANSFER_LIMITS uint64 = 20000000
	TIER_ONE_SINGLE_TRANSFER_LIMITS uint64 = 5000000

	// verify phone number, add address and NIN
	TIER_TWO_DAILY_TRANSFER_LIMITS uint64 = 90000000
	TIER_TWO_SINGLE_TRANSFER_LIMITS uint64 = 90000000

	// address verified
	TIER_THREE_DAILY_TRANSFER_LIMITS uint64 = 100000000
	TIER_THREE_SINGLE_TRANSFER_LIMITS uint64 = 2500000000
)

// polymer response codes
// these consist of 4 digitsnumbers
//
// the 1st 3 are randomly generated but represent specific scenarios
// 4th indicates if the response requires user interactions through a dialog box. 0 means it does not require. 1 means it requires.

var REQUIRES_FACE_MATCH_UNLOCK uint = 9870 // take the user to the face match page to unlock the account
var ENCRYPTION_KEY_EXPIRED uint = 6170 // take the user to the face match page to unlock the account
var NIN_VERIFICATION_FAILED uint = 6511 // show the user the popup to escale the issue for manual review
var EMAIL_UNVERIFIED uint = 3540 // request a new otp and redirect the user to the otp verification screen

type state struct {
	Name     string   `json:"state"`
	LGAs     []string `json:"lgas"`
}

var States = []state{
	{
		Name: "Lagos",
		LGAs: []string{
			"Agege",
			"Ajeromi-Ifelodun",
			"Alimosho",
			"Amuwo-Odofin",
			"Badagry",
			"Apapa",
			"Epe",
			"Eti Osa",
			"Ibeju-Lekki",
			"Ifako-Ijaiye",
			"Ikeja",
			"Ikorodu",
			"Kosofe",
			"Lagos Island",
			"Mushin",
			"Lagos Mainland",
			"Ojo",
			"Oshodi-Isolo",
			"Shomolu",
			"Surulere",
		},
	},
	{
		Name: "Abuja",
		LGAs: []string{
			"Abaji",
			"Bwari",
			"Gwagwalada",
			"Kuje",
			"Kwali",
			"Municipal Area Council",
		},
	},
	{
		Name: "Abia",
		LGAs: []string{
			"Aba North",
			"Arochukwu",
			"Aba South",
			"Bende",
			"Isiala Ngwa North",
			"Ikwuano",
			"Isiala Ngwa South",
			"Isuikwuato",
			"Obi Ngwa",
			"Ohafia",
			"Osisioma",
			"Ugwunagbo",
			"Ukwa East",
			"Ukwa West",
			"Umuahia North",
			"Umuahia South",
			"Umu Nneochi",
		},
	},
	{
		Name: "Adamawa",
		LGAs: []string{
			"Demsa",
			"Fufure",
			"Ganye",
			"Gayuk",
			"Gombi",
			"Grie",
			"Hong",
			"Jada",
			"Larmurde",
			"Madagali",
			"Maiha",
			"Mayo Belwa",
			"Michika",
			"Mubi North",
			"Mubi South",
			"Numan",
			"Shelleng",
			"Song",
			"Toungo",
			"Yola North",
			"Yola South",
		},
	},
	{
		Name: "Akwa Ibom",
		LGAs: []string{
			"Abak",
			"Eastern Obolo",
			"Eket",
			"Esit Eket",
			"Essien Udim",
			"Etim Ekpo",
			"Etinan",
			"Ibeno",
			"Ibesikpo Asutan",
			"Ibiono-Ibom",
			"Ikot Abasi",
			"Ika",
			"Ikono",
			"Ikot Ekpene",
			"Ini",
			"Mkpat-Enin",
			"Itu",
			"Mbo",
			"Nsit-Atai",
			"Nsit-Ibom",
			"Nsit-Ubium",
			"Obot Akara",
			"Okobo",
			"Onna",
			"Oron",
			"Udung-Uko",
			"Ukanafun",
			"Oruk Anam",
			"Uruan",
			"Urue-Offong/Oruko",
			"Uyo",
		},
	},
	{
		Name: "Anambra",
		LGAs: []string{
			"Aguata",
			"Anambra East",
			"Anaocha",
			"Awka North",
			"Anambra West",
			"Awka South",
			"Ayamelum",
			"Dunukofia",
			"Ekwusigo",
			"Idemili North",
			"Idemili South",
			"Ihiala",
			"Njikoka",
			"Nnewi North",
			"Nnewi South",
			"Ogbaru",
			"Onitsha North",
			"Onitsha South",
			"Orumba North",
			"Orumba South",
			"Oyi",
		},
	},
	{
		Name: "Nasarawa",
		LGAs: []string{
			"Akwanga",
			"Awe",
			"Doma",
			"Karu",
			"Keana",
			"Keffi",
			"Lafia",
			"Kokona",
			"Nasarawa Egon",
			"Nasarawa",
			"Obi",
			"Toto",
			"Wamba",
		},
	},
	{
		Name: "Ogun",
		LGAs: []string{
			"Abeokuta North",
			"Abeokuta South",
			"Ado-Odo/Ota",
			"Egbado North",
			"Ewekoro",
			"Egbado South",
			"Ijebu North",
			"Ijebu East",
			"Ifo",
			"Ijebu Ode",
			"Ijebu North East",
			"Imeko Afon",
			"Ikenne",
			"Ipokia",
			"Odeda",
			"Obafemi Owode",
			"Odogbolu",
			"Remo North",
			"Ogun Waterside",
			"Shagamu",
		},
	},
	{
		Name: "Ondo",
		LGAs: []string{
			"Akoko North-East",
			"Akoko North-West",
			"Akoko South-West",
			"Akoko South-East",
			"Akure North",
			"Akure South",
			"Ese Odo",
			"Idanre",
			"Ifedore",
			"Ilaje",
			"Irele",
			"Ile Oluji/Okeigbo",
			"Odigbo",
			"Okitipupa",
			"Ondo West",
			"Ose",
			"Ondo East",
			"Owo",
		},
	},
	{
		Name: "Rivers",
		LGAs: []string{
			"Abua/Odual",
			"Ahoada East",
			"Ahoada West",
			"Andoni",
			"Akuku-Toru",
			"Asari-Toru",
			"Bonny",
			"Degema",
			"Emuoha",
			"Eleme",
			"Ikwerre",
			"Etche",
			"Gokana",
			"Khana",
			"Obio/Akpor",
			"Ogba/Egbema/Ndoni",
			"Ogu/Bolo",
			"Okrika",
			"Omuma",
			"Opobo/Nkoro",
			"Oyigbo",
			"Port Harcourt",
			"Tai",
		},
	},
	{
		Name: "Bauchi",
		LGAs: []string{
			"Alkaleri",
			"Bauchi",
			"Bogoro",
			"Damban",
			"Darazo",
			"Dass",
			"Gamawa",
			"Ganjuwa",
			"Giade",
			"Itas/Gadau",
			"Jama'are",
			"Katagum",
			"Kirfi",
			"Misau",
			"Ningi",
			"Shira",
			"Tafawa Balewa",
			"Toro",
			"Warji",
			"Zaki",
		},
	},
	{
		Name: "Benue",
		LGAs: []string{
			"Agatu",
			"Apa",
			"Ado",
			"Buruku",
			"Gboko",
			"Guma",
			"Gwer East",
			"Gwer West",
			"Katsina-Ala",
			"Konshisha",
			"Kwande",
			"Logo",
			"Makurdi",
			"Obi",
			"Ogbadibo",
			"Ohimini",
			"Oju",
			"Okpokwu",
			"Oturkpo",
			"Tarka",
			"Ukum",
			"Ushongo",
			"Vandeikya",
		},
	},
	{
		Name: "Borno",
		LGAs: []string{
			"Abadam",
			"Askira/Uba",
			"Bama",
			"Bayo",
			"Biu",
			"Chibok",
			"Damboa",
			"Dikwa",
			"Guzamala",
			"Gubio",
			"Hawul",
			"Gwoza",
			"Jere",
			"Kaga",
			"Kala/Balge",
			"Konduga",
			"Kukawa",
			"Kwaya Kusar",
			"Mafa",
			"Magumeri",
			"Maiduguri",
			"Mobbar",
			"Marte",
			"Monguno",
			"Ngala",
			"Nganzai",
			"Shani",
		},
	},
	{
		Name: "Bayelsa",
		LGAs: []string{
			"Brass",
			"Ekeremor",
			"Kolokuma/Opokuma",
			"Nembe",
			"Ogbia",
			"Sagbama",
			"Southern Ijaw",
			"Yenagoa",
		},
	},
	{
		Name: "Cross River",
		LGAs: []string{
			"Abi",
			"Akamkpa",
			"Akpabuyo",
			"Bakassi",
			"Bekwarra",
			"Biase",
			"Boki",
			"Calabar Municipal",
			"Calabar South",
			"Etung",
			"Ikom",
			"Obanliku",
			"Obubra",
			"Obudu",
			"Odukpani",
			"Ogoja",
			"Yakuur",
			"Yala",
		},
	},
	{
		Name: "Delta",
		LGAs: []string{
			"Aniocha North",
			"Aniocha South",
			"Bomadi",
			"Burutu",
			"Ethiope West",
			"Ethiope East",
			"Ika North East",
			"Ika South",
			"Isoko North",
			"Isoko South",
			"Ndokwa East",
			"Ndokwa West",
			"Okpe",
			"Oshimili North",
			"Oshimili South",
			"Patani",
			"Sapele",
			"Udu",
			"Ughelli North",
			"Ukwuani",
			"Ughelli South",
			"Uvwie",
			"Warri North",
			"Warri South",
			"Warri South West",
		},
	},
	{
		Name: "Ebonyi",
		LGAs: []string{
			"Abakaliki",
			"Afikpo North",
			"Ebonyi",
			"Afikpo South",
			"Ezza North",
			"Ikwo",
			"Ezza South",
			"Ivo",
			"Ishielu",
			"Izzi",
			"Ohaozara",
			"Ohaukwu",
			"Onicha",
		},
	},
	{
		Name: "Edo",
		LGAs: []string{
			"Akoko-Edo",
			"Egor",
			"Esan Central",
			"Esan North-East",
			"Esan South-East",
			"Esan West",
			"Etsako Central",
			"Etsako East",
			"Etsako West",
			"Igueben",
			"Ikpoba Okha",
			"Orhionmwon",
			"Oredo",
			"Ovia North-East",
			"Ovia South-West",
			"Owan East",
			"Owan West",
			"Uhunmwonde",
		},
	},
	{
		Name: "Ekiti",
		LGAs: []string{
			"Ado Ekiti",
			"Efon",
			"Ekiti East",
			"Ekiti South-West",
			"Ekiti West",
			"Emure",
			"Gbonyin",
			"Ido Osi",
			"Ijero",
			"Ikere",
			"Ilejemeje",
			"Irepodun/Ifelodun",
			"Ikole",
			"Ise/Orun",
			"Moba",
			"Oye",
		},
	},
	{
		Name: "Enugu",
		LGAs: []string{
			"Awgu",
			"Aninri",
			"Enugu East",
			"Enugu North",
			"Ezeagu",
			"Enugu South",
			"Igbo Etiti",
			"Igbo Eze North",
			"Igbo Eze South",
			"Isi Uzo",
			"Nkanu East",
			"Nkanu West",
			"Nsukka",
			"Udenu",
			"Oji River",
			"Uzo Uwani",
			"Udi",
		},
	},
	{
		Name: "Gombe",
		LGAs: []string{
			"Akko",
			"Balanga",
			"Billiri",
			"Dukku",
			"Funakaye",
			"Gombe",
			"Kaltungo",
			"Kwami",
			"Nafada",
			"Shongom",
			"Yamaltu",
			"Deba",
		},
	},
	{
		Name: "Jigawa",
		LGAs: []string{
			"Auyo",
			"Babura",
			"Buji",
			"Biriniwa",
			"Birnin Kudu",
			"Dutse",
			"Gagarawa",
			"Garki",
			"Gumel",
			"Guri",
			"Gwaram",
			"Gwiwa",
			"Hadejia",
			"Jahun",
			"Kafin Hausa",
			"Kazaure",
			"Kiri Kasama",
			"Kiyawa",
			"Kaugama",
			"Maigatari",
			"Malam Madori",
			"Miga",
			"Sule Tankarkar",
			"Roni",
			"Ringim",
			"Yankwashi",
			"Taura",
		},
	},
	{
		Name: "Oyo",
		LGAs: []string{
			"Afijio",
			"Akinyele",
			"Atiba",
			"Atisbo",
			"Egbeda",
			"Ibadan North",
			"Ibadan North-East",
			"Ibadan North-West",
			"Ibadan South-East",
			"Ibarapa Central",
			"Ibadan South-West",
			"Ibarapa East",
			"Ido",
			"Ibarapa North",
			"Irepo",
			"Iseyin",
			"Itesiwaju",
			"Iwajowa",
			"Kajola",
			"Lagelu",
			"Ogbomosho North",
			"Ogbomosho South",
			"Ogo Oluwa",
			"Olorunsogo",
			"Oluyole",
			"Ona Ara",
			"Orelope",
			"Ori Ire",
			"Oyo",
			"Oyo East",
			"Saki East",
			"Saki West",
			"Surulere Oyo State",
		},
	},
	{
		Name: "Imo",
		LGAs: []string{
			"Aboh Mbaise",
			"Ahiazu Mbaise",
			"Ehime Mbano",
			"Ezinihitte",
			"Ideato North",
			"Ideato South",
			"Ihitte/Uboma",
			"Ikeduru",
			"Isiala Mbano",
			"Mbaitoli",
			"Isu",
			"Ngor Okpala",
			"Njaba",
			"Nkwerre",
			"Nwangele",
			"Obowo",
			"Oguta",
			"Ohaji/Egbema",
			"Okigwe",
			"Orlu",
			"Orsu",
			"Oru East",
			"Oru West",
			"Owerri Municipal",
			"Owerri North",
			"Unuimo",
			"Owerri West",
		},
	},
	{
		Name: "Kaduna",
		LGAs: []string{
			"Birnin Gwari",
			"Chikun",
			"Giwa",
			"Ikara",
			"Igabi",
			"Jaba",
			"Jema'a",
			"Kachia",
			"Kaduna North",
			"Kaduna South",
			"Kagarko",
			"Kajuru",
			"Kaura",
			"Kauru",
			"Kubau",
			"Kudan",
			"Lere",
			"Makarfi",
			"Sabon Gari",
			"Sanga",
			"Soba",
			"Zangon Kataf",
			"Zaria",
		},
	},
	{
		Name: "Kebbi",
		LGAs: []string{
			"Aleiro",
			"Argungu",
			"Arewa Dandi",
			"Augie",
			"Bagudo",
			"Birnin Kebbi",
			"Bunza",
			"Dandi",
			"Fakai",
			"Gwandu",
			"Jega",
			"Kalgo",
			"Koko/Besse",
			"Maiyama",
			"Ngaski",
			"Shanga",
			"Suru",
			"Sakaba",
			"Wasagu/Danko",
			"Yauri",
			"Zuru",
		},
	},
	{
		Name: "Kano",
		LGAs: []string{
			"Ajingi",
			"Albasu",
			"Bagwai",
			"Bebeji",
			"Bichi",
			"Bunkure",
			"Dala",
			"Dambatta",
			"Dawakin Kudu",
			"Dawakin Tofa",
			"Doguwa",
			"Fagge",
			"Gabasawa",
			"Garko",
			"Garun Mallam",
			"Gezawa",
			"Gaya",
			"Gwale",
			"Gwarzo",
			"Kabo",
			"Kano Municipal",
			"Karaye",
			"Kibiya",
			"Kiru",
			"Kumbotso",
			"Kunchi",
			"Kura",
			"Madobi",
			"Makoda",
			"Minjibir",
			"Nasarawa",
			"Rano",
			"Rimin Gado",
			"Rogo",
			"Shanono",
			"Takai",
			"Sumaila",
			"Tarauni",
			"Tofa",
			"Tsanyawa",
			"Tudun Wada",
			"Ungogo",
			"Warawa",
			"Wudil",
		},
	},
	{
		Name: "Katsina",
		LGAs: []string{
			"Bakori",
			"Batagarawa",
			"Batsari",
			"Baure",
			"Bindawa",
			"Charanchi",
			"Danja",
			"Dandume",
			"Dan Musa",
			"Daura",
			"Dutsi",
			"Dutsin Ma",
			"Faskari",
			"Funtua",
			"Ingawa",
			"Jibia",
			"Kafur",
			"Kaita",
			"Kankara",
			"Kankia",
			"Katsina",
			"Kurfi",
			"Kusada",
			"Mai'Adua",
			"Malumfashi",
			"Mani",
			"Mashi",
			"Matazu",
			"Musawa",
			"Rimi",
			"Sabuwa",
			"Safana",
			"Sandamu",
			"Zango",
		},
	},
	{
		Name: "Kwara",
		LGAs: []string{
			"Asa",
			"Baruten",
			"Edu",
			"Ilorin East",
			"Ifelodun",
			"Ilorin South",
			"Ekiti Kwara State",
			"Ilorin West",
			"Irepodun",
			"Isin",
			"Kaiama",
			"Moro",
			"Offa",
			"Oke Ero",
			"Oyun",
			"Pategi",
		},
	},
	{
		Name: "Kogi",
		LGAs: []string{
			"Ajaokuta",
			"Adavi",
			"Ankpa",
			"Bassa",
			"Dekina",
			"Ibaji",
			"Idah",
			"Igalamela Odolu",
			"Ijumu",
			"Kogi",
			"Kabba/Bunu",
			"Lokoja",
			"Ofu",
			"Mopa Muro",
			"Ogori/Magongo",
			"Okehi",
			"Okene",
			"Olamaboro",
			"Omala",
			"Yagba East",
			"Yagba West",
		},
	},
	{
		Name: "Osun",
		LGAs: []string{
			"Aiyedire",
			"Atakunmosa West",
			"Atakunmosa East",
			"Aiyedaade",
			"Boluwaduro",
			"Boripe",
			"Ife East",
			"Ede South",
			"Ife North",
			"Ede North",
			"Ife South",
			"Ejigbo",
			"Ife Central",
			"Ifedayo",
			"Egbedore",
			"Ila",
			"Ifelodun",
			"Ilesa East",
			"Ilesa West",
			"Irepodun",
			"Irewole",
			"Isokan",
			"Iwo",
			"Obokun",
			"Odo Otin",
			"Ola Oluwa",
			"Olorunda",
			"Oriade",
			"Orolu",
			"Osogbo",
		},
	},
	{
		Name: "Sokoto",
		LGAs: []string{
			"Gudu",
			"Gwadabawa",
			"Illela",
			"Isa",
			"Kebbe",
			"Kware",
			"Rabah",
			"Sabon Birni",
			"Shagari",
			"Silame",
			"Sokoto North",
			"Sokoto South",
			"Tambuwal",
			"Tangaza",
			"Tureta",
			"Wamako",
			"Wurno",
			"Yabo",
			"Binji",
			"Bodinga",
			"Dange Shuni",
			"Goronyo",
			"Gada",
		},
	},
	{
		Name: "Plateau",
		LGAs: []string{
			"Bokkos",
			"Barkin Ladi",
			"Bassa",
			"Jos East",
			"Jos North",
			"Jos South",
			"Kanam",
			"Kanke",
			"Langtang South",
			"Langtang North",
			"Mangu",
			"Mikang",
			"Pankshin",
			"Qua'an Pan",
			"Riyom",
			"Shendam",
			"Wase",
		},
	},
	{
		Name: "Taraba",
		LGAs: []string{
			"Ardo Kola",
			"Bali",
			"Donga",
			"Gashaka",
			"Gassol",
			"Ibi",
			"Jalingo",
			"Karim Lamido",
			"Kumi",
			"Lau",
			"Sardauna",
			"Takum",
			"Ussa",
			"Wukari",
			"Yorro",
			"Zing",
		},
	},
	{
		Name: "Yobe",
		LGAs: []string{
			"Bade",
			"Bursari",
			"Damaturu",
			"Fika",
			"Fune",
			"Geidam",
			"Gujba",
			"Gulani",
			"Jakusko",
			"Karasuwa",
			"Machina",
			"Nangere",
			"Nguru",
			"Potiskum",
			"Tarmuwa",
			"Yunusari",
			"Yusufari",
		},
	},
	{
		Name: "Zamfara",
		LGAs: []string{
			"Anka",
			"Birnin Magaji/Kiyaw",
			"Bakura",
			"Bukkuyum",
			"Bungudu",
			"Gummi",
			"Gusau",
			"Kaura Namoda",
			"Maradun",
			"Shinkafi",
			"Maru",
			"Talata Mafara",
			"Tsafe",
			"Zurmi",
		},
	},
}
  