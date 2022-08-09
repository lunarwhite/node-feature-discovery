package usb

import "testing"

func TestGetReadableVendor(t *testing.T) {
	var tests = []struct {
		input  string
		wanted string
	}{
		{"ASUSTek Computer, Inc. (wrong ID)", "Wrong-ID"},
		{"Nebraska Furniture Mart", "Nebraska-F"},
		{"Toshiba Corp., Digital Media Equipment", "Toshiba"},
		{"Connector Co., Ltd", "Connector"},
		{"Chipsbank Microelectronics Co., Ltd", "Chipsbank"},
		{"Trust International B.V.", "Trust-Inte"},
		{"Bernd Walter Computer Technology", "Bernd-Walt"},
		{"RFC Distribution(s) PTE, Ltd", "RFC-Distri"},
		{"EndPoints, Inc.", "EndPoints"},
		{"Samsung Info. Systems America, Inc.", "Samsung-In"},
		{"AT&T Paradyne", "AT&T-Parad"},
		{"AMP/Tycoelectronics Corp.", "AMP-Tycoel"},
		{"Foxconn / Hon Hai", "Foxconn"},
		{"Capet (Kaohsiung) Corp.", "Capet-Kaoh"},
		{"Acer Peripherals Inc. (now BenQ Corp.)", "Acer-Perip"},
		{"Temic MHS S.A.", "Temic-MHS"},
		{"Hand (Welch Allyn, Inc.)", "Hand-Welch"},
		{"Kodak, Ltd.", "Kodak"},
		{"Y.C. Cable U.S.A., Inc.", "YC-Cable"},
		{"U.S. Robotics (3Com)", "US-Roboti"},
		{"Who? Vision Systems, Inc.", "Who-Vision"},
		{"TSAY-E (BVI) International, Inc.", "TSAY-E-BVI"},
		{"Lernout + Hauspie", "Lernout"},
		{"Yubico.com", "Yubicocom"},
		{"T+A elektroakustik GmbH & Co KG, Germany", "T+A-elektr"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			output := getReadableVendor(tt.input)
			if output != tt.wanted {
				t.Errorf("got: %s, wanted: %s", output, tt.wanted)
			}
		})
	}
}

func TestGetReadableDevice(t *testing.T) {
	var tests = []struct {
		input  string
		wanted string
	}{
		{"DataTraveler 2.0 1GB/4GB Flash Drive / Patriot Xporter 4GB Flash Drive", "DataTraveler-2.0-1GB-4GB-Flash-D"},
		{"Flash Drive 2 GB [ICIDU 2 GB]", "Flash-Drive-2-GB-ICIDU-2-GB"},
		{"USB flash drive (32 GB SHARKOON Accelerate)", "USB-flash-drive-32-GB-SHARKOON-A"},
		{"flash drive (2GB, EMTEC)", "flash-drive-2GB-EMTEC"},
		{"Mass-Storage Device [NT2 U3.1]", "Mass-Storage-Device-NT2-U3.1"},
		{"AnyPoint (TM) Home Network 1.6 Mbps Wireless Adapter", "AnyPoint-TM-Home-Network-1.6-Mbp"},
		{"3.0 root hub", "3.0-root-hub"},
		{"USB 2.0 Hub", "USB-2.0-Hub"},
		{"Myriad VPU [Movidius Neural Compute Stick]", "Myriad-VPU-Movidius-Neural-Compu"},
		{"AnyPoint(TM) Wireless II Network 11Mbps Adapter [Atmel AT76C503A]", "AnyPoint-TM-Wireless-II-Network"},
		{"Bluetooth 4.0* Smart Ready (low energy)", "Bluetooth-4.0-Smart-Ready-low-en"},
		{"AnyPoint® 3240 Modem - WAN", "AnyPoint-3240-Modem-WAN"},
		{"8 Series/C220 Series EHCI #1", "8-Series-C220-Series-EHCI-1"},
		{"SideWinder® Freestyle Pro", "SideWinder-Freestyle-Pro"},
		{"Xbox & PC Gamepad", "Xbox-PC-Gamepad"},
		{"H8314 [Xperia XZ2 Compact] (MIDI)", "H8314-Xperia-XZ2-Compact-MIDI"},
		{"LP1965 19\" Monitor Hub", "LP1965-19-Monitor-Hub"},
		{"Mouse*in*a*Box Optical Pro", "Mouse-in-a-Box-Optical-Pro"},
		{"PS/2 Keyboard, Mouse & Joystick Ports", "PS-2-Keyboard-Mouse-Joystick-Por"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			output := getReadableDevice(tt.input)
			if output != tt.wanted {
				t.Errorf("got: %s, wanted: %s", output, tt.wanted)
			}
		})
	}
}

func TestInArray(t *testing.T) {
	var hiddenSuffix = []string{
		"Corporation", "Corp.", "Corp.,", "Corp", "corp.", "Co.", "Co", "co.", "co.,", "CO.,", "Co.,", "Co.,Ltd", "Co.,LTD.", "Co.,Ltd.",
		"INC.", "INC", "Inc.", "Inc", "Inc,", "inc.",
		"Ltd.", "Ltd", "LTD.", "ltd.",
		"Technologies", "Technologies,", "Technology", "Technology,",
		"Information",
		"Electronics", "ELECTRONICS ", "Electric", "ELECTRIC",
		"Company",
		"Group",
		"LLC", "LLC.",
	}
	var tests = []struct {
		input  string
		wanted bool
	}{
		{"Corp", true},
		{"Corporation", true},
		{"Co.,Ltd", true},
		{"Inc", true},
		{"Technologies", true},
		{"LLC", true},
		{"Scoop", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			output := inArray(tt.input, hiddenSuffix)
			if output != tt.wanted {
				t.Errorf("got: %t, wanted: %t", output, tt.wanted)
			}
		})
	}
}
