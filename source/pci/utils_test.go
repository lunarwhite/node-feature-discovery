package pci

import "testing"

func TestGetReadableClass(t *testing.T) {
	var tests = []struct {
		input  string
		wanted string
	}{
		{"Ethernet controller", "Ethe"},
		{"VGA compatible controller", "VGA"},
		{"3D controller", "3D"},
		{"Token ring network controller", "Toke"},
		{"Network and computing encryption device", "Netw"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			output := getReadableClass(tt.input)
			if output != tt.wanted {
				t.Errorf("got: %s, wanted: %s", output, tt.wanted)
			}
		})
	}
}

func TestGetReadableVendor(t *testing.T) {
	var tests = []struct {
		input  string
		wanted string
	}{
		{"Silicon Image, Inc. (Wrong ID)", "Wrong-ID"},        // wrong-id
		{"VMWare Inc (temporary ID)", "Wrong-ID"},             // wrong-id
		{"Advanced Micro Devices, Inc. [AMD/ATI]", "AMD-ATI"}, // []
		{"T1042 [Freescale]", "Freescale"},                    // []
		{"Sapphire, Inc.", "Sapphire"},
		{"Matsuta-Kotobuki Electronics Industries, Ltd.", "Matsuta-Ko"},
		{"Tata Power Strategic Electronics Division", "Tata-Power"},
		{"Synopsys/Logic Modeling Group", "Synopsys-L"},
		{"Synopsysi/Logic Modeling Group", "Synopsysi"},
		{"Wired Inc.", "Wired"},
		{"AT&T GIS (NCR)", "AT&T-GIS"},
		{"IPC Corporation Ltd.", "IPC"},
		{"NexGen Microsystems", "NexGen-Mic"},
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
		{"BeaverCreek HDMI Audio [Radeon HD 6500D and 6400G-6600G series]", "Radeon-HD-6500D-and-6400G-6600G-se"},
		{"IXP SB4x0 High Definition Audio Controller", "IXP-SB4x0-High-Definition-Audio"},
		{"Compute Engine Virtual Ethernet [gVNIC]", "gVNIC"},
		{"Airbrush Combined Paintbox IPU/Oscar Edge TPU [Pixel Neural Core]", "Pixel-Neural-Core"},
		{"Ethernet Controller X710 Intel(R) FPGA Programmable Acceleration Card N3000 for Networking", "Ethernet-X710-IntelR-FPGA-Programm"},
		{"NV5 [Riva TNT2 Model 64 / Model 64 Pro]", "Riva-TNT2-Model-64-Model-64-Pro"},
		{"Atom Processor Z36xxx/Z37xxx Series SDIO Controller", "Atom-Z36xxx-Z37xxx-Series-SDIO"},
		{"7500/5520/5500 Routing & Protocol Layer Register Port 1", "7500-5520-5500-Routing-Protocol-La"},
		{"450NX - 82451NX Memory & I/O Controller", "450NX-82451NX-Memory-I-O"},
		{"Iris Plus Graphics G7 (Ice Lake)", "Iris-Plus-G7-Ice-Lake"},
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
		"Corp.", "Corp", "Corporation", "Co",
		"Inc.", "Inc",
		"Co,.Ltd",
		"Ltd.", "Ltd",
		"Technologies", "Technology",
		"Information",
		"Company",
		"Group",
		"LLC",
	}
	var tests = []struct {
		input  string
		wanted bool
	}{
		{"Corp", true},
		{"Corporation", true},
		{"Co", true},
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
