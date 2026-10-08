package chain

import "testing"

func TestNetworkPresets(t *testing.T) {
	t.Setenv("BODHI_CHAIN_NETWORK", "polygon")
	c := ConfigFromEnv()
	if c.ChainID != 137 || c.ExplorerURL != "https://polygonscan.com" || len(c.RPCs()) != 4 {
		t.Errorf("polygon preset = %+v", c)
	}
	t.Setenv("BODHI_CHAIN_NETWORK", "")
	t.Setenv("BODHI_CHAIN_RPC", "https://a.example, https://b.example")
	c = ConfigFromEnv()
	if c.ChainID != 80002 || c.ExplorerURL != "https://amoy.polygonscan.com" || len(c.RPCs()) != 2 {
		t.Errorf("default preset with own nodes = %+v", c)
	}
}
