package core

type ProtoVersion struct {
	Prov *ProtoVersionProv `json:"prov"`
}

type ProtoVersionProv struct {
	// Pointers are used as nullable types
	Ver         *string  `json:"ver,omitempty"`
	SecVer      *int     `json:"secVer,omitempty"`
	Cap         []string `json:"cap,omitempty"`
	SecPatchVer *int     `json:"secPatchVer,omitempty"`
}
