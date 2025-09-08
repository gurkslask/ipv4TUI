package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
	ColorR1  = lipgloss.Color("#DC243C")
	ColorR2 = lipgloss.Color("#F75270")
	ColorR3 = lipgloss.Color("#F7CAC9")
	ColorR4 = lipgloss.Color("#FDEBD0")
	ColorBG1 = lipgloss.Color("#BADA55")
	ColorBG2 = lipgloss.Color("#FFFA55")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray).Background(ColorBG2)
	oct1Style     = lipgloss.NewStyle().Foreground(ColorR1)
	oct2Style     = lipgloss.NewStyle().Foreground(ColorR2)
	oct3Style     = lipgloss.NewStyle().Foreground(ColorR3)
	oct4Style     = lipgloss.NewStyle().Foreground(ColorR4)
	staticStyle    = lipgloss.NewStyle()
)

type model struct {
	inputs       []textinput.Model // ipv4 Octets
	focused      int               // which item our cursor is pointing at
	err          error             // error?
	ipaddr       IPv4              // the ip address
	ipBinary     string
	subnetmask   IPv4
	subnetBinary string
	netaddr	IPv4
	broadcastaddr IPv4
}

func initialModel() model {
	ip, _ := NewIPv4("1.2.3.4")
	snm, _ := NewIPv4("255.255.255.0")
	var inputs []textinput.Model = make([]textinput.Model, 8)
	ipBinary := ip.PrintBinary()
	subnetBinary := snm.PrintBinary()

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "1"
	inputs[0].SetValue("192")
	inputs[0].Focus()
	inputs[0].CharLimit = 3
	inputs[0].Width = 7
	inputs[0].Prompt = ""
	inputs[0].Validate = octetValidator
	inputs[0].TextStyle = inputStyle

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "2"
	inputs[1].SetValue("168")
	inputs[1].CharLimit = 3
	inputs[1].Width = 7
	inputs[1].Prompt = ""
	inputs[1].Validate = octetValidator
	inputs[1].TextStyle = inputStyle

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "3"
	inputs[2].SetValue("0")
	inputs[2].CharLimit = 3
	inputs[2].Width = 7
	inputs[2].Prompt = ""
	inputs[2].Validate = octetValidator
	inputs[2].TextStyle = inputStyle

	inputs[3] = textinput.New()
	inputs[3].Placeholder = "4"
	inputs[3].SetValue("2")
	inputs[3].CharLimit = 3
	inputs[3].Width = 11
	inputs[3].Prompt = ""
	inputs[3].Validate = octetValidator
	inputs[3].TextStyle = inputStyle

	inputs[4] = textinput.New()
	inputs[4].Placeholder = "4"
	inputs[4].SetValue("255")
	inputs[4].CharLimit = 3
	inputs[4].Width = 7
	inputs[4].Prompt = ""
	inputs[4].Validate = octetValidator
	inputs[4].TextStyle = inputStyle

	inputs[5] = textinput.New()
	inputs[5].Placeholder = "4"
	inputs[5].SetValue("255")
	inputs[5].CharLimit = 3
	inputs[5].Width = 7
	inputs[5].Prompt = ""
	inputs[5].Validate = octetValidator
	inputs[5].TextStyle = inputStyle

	inputs[6] = textinput.New()
	inputs[6].Placeholder = "4"
	inputs[6].SetValue("255")
	inputs[6].CharLimit = 3
	inputs[6].Width = 7
	inputs[6].Prompt = ""
	inputs[6].Validate = octetValidator
	inputs[6].TextStyle = inputStyle

	inputs[7] = textinput.New()
	inputs[7].Placeholder = "4"
	inputs[7].SetValue("0")
	inputs[7].CharLimit = 3
	inputs[7].Width = 7
	inputs[7].Prompt = ""
	inputs[7].Validate = octetValidator
	inputs[7].TextStyle = inputStyle
	return model{
		ipaddr:     ip,
		subnetmask: snm,
		inputs:     inputs,
		focused:    0,
		err:        nil,
		ipBinary: ipBinary,
		subnetBinary: subnetBinary,
	}
}
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			//if m.focused == len(m.inputs)-1 {
			//return m, tea.Quit
			//}
			m.updateIPAddress()
			m.updateIPSubnetmask()
			m.netaddr = CalcNetAddress(m.ipaddr, m.subnetmask)
			m.broadcastaddr = CalcBroadcastAddress(m.ipaddr, m.subnetmask)
			m.ipBinary = m.ipaddr.PrintBinary()
			m.subnetBinary = m.subnetmask.PrintBinary()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()
	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}
func (m model) View() string {
	return fmt.Sprintf(
		` IP-Calculator

 %s%s
 %s%s%s%s%s%s%s%s
 %s %s %s %s  %s %s %s %s

 %s%s
 %s.%s.%s.%s  %s.%s.%s.%s

 Binary
 netaddr
 %s
 broadcast
 %s
 %s
 %s
`,

		staticStyle.Width(41).Render("IP-address"),
		staticStyle.Width(20).Render("Subnetmask"),
		staticStyle.Width(5).Render("Oct1"),
		staticStyle.Width(5).Render("Oct2"),
		staticStyle.Width(5).Render("Oct3"),
		staticStyle.Width(26).Render("Oct4"),
		staticStyle.Width(5).Render("Oct1"),
		staticStyle.Width(5).Render("Oct2"),
		staticStyle.Width(5).Render("Oct3"),
		staticStyle.Width(5).Render("Oct4"),
		m.inputs[0].View(),
		m.inputs[1].View(),
		m.inputs[2].View(),
		m.inputs[3].View(),
		m.inputs[4].View(),
		m.inputs[5].View(),
		m.inputs[6].View(),
		m.inputs[7].View(),
		staticStyle.Width(41).Render("IP-address"),
		staticStyle.Width(20).Render("Subnetmask"),
		oct1Style.Width(8).Render(string(m.ipBinary[0:8])),
		oct2Style.Width(8).Render(string(m.ipBinary[8:16])),
		oct3Style.Width(8).Render(string(m.ipBinary[16:24])),
		oct4Style.Width(12).Render(string(m.ipBinary[24:32])),
		oct1Style.Width(8).Render(string(m.subnetBinary[0:8])),
		oct2Style.Width(8).Render(string(m.subnetBinary[8:16])),
		oct3Style.Width(8).Render(string(m.subnetBinary[16:24])),
		oct4Style.Width(8).Render(string(m.subnetBinary[24:32])),
		m.netaddr.PrintDecimal(),
		m.broadcastaddr.PrintDecimal(),
		m.err,
		continueStyle.Render("Continue ->"),
	) + "\n"
}

// nextInput focuses the next input field
func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *model) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
func (m *model) updateIPAddress() {
	ss := fmt.Sprintf("%v.%v.%v.%v", m.inputs[0].Value(), m.inputs[1].Value(), m.inputs[2].Value(), m.inputs[3].Value())
	m.err = m.ipaddr.CheckIfValidIPAddress(ss)
	if m.err != nil {
		return
	}
	m.err = m.ipaddr.SetAddress(ss)
}
func (m *model) updateIPSubnetmask() {
	ss := fmt.Sprintf("%v.%v.%v.%v", m.inputs[4].Value(), m.inputs[5].Value(), m.inputs[6].Value(), m.inputs[7].Value())
	m.err = m.subnetmask.SetAddress(ss)
	m.err = m.subnetmask.CheckIfValidSubnetmask()
}
/*func (m model) View() string {
	return fmt.Sprintf(
		` IP-Calculator

IP-address
 %s %s %s %s
 %s %s %s %s

 Binary
 %s.%s.%s.%s

 Subnetmask
 %s %s %s %s
 %s %s %s %s

 Binary
 %s.%s.%s.%s

 netaddr
 %s
 broadcast
 %s
 %s
 %s
`,
		inputStyle.Width(5).Render("Oct1"),
		inputStyle.Width(5).Render("Oct2"),
		inputStyle.Width(5).Render("Oct3"),
		inputStyle.Width(5).Render("Oct4"),
		m.inputs[0].View(),
		m.inputs[1].View(),
		m.inputs[2].View(),
		m.inputs[3].View(),
		oct1Style.Width(9).Render(string(m.ipBinary[0:8])),
		oct2Style.Width(9).Render(string(m.ipBinary[8:16])),
		oct3Style.Width(9).Render(string(m.ipBinary[16:24])),
		oct4Style.Width(9).Render(string(m.ipBinary[24:32])),
		inputStyle.Width(5).Render("Oct1"),
		inputStyle.Width(5).Render("Oct2"),
		inputStyle.Width(5).Render("Oct3"),
		inputStyle.Width(5).Render("Oct4"),
		m.inputs[4].View(),
		m.inputs[5].View(),
		m.inputs[6].View(),
		m.inputs[7].View(),
		oct1Style.Width(9).Render(string(m.subnetBinary[0:8])),
		oct2Style.Width(9).Render(string(m.subnetBinary[8:16])),
		oct3Style.Width(9).Render(string(m.subnetBinary[16:24])),
		oct4Style.Width(9).Render(string(m.subnetBinary[24:32])),
		m.netaddr.PrintDecimal(),
		m.broadcastaddr.PrintDecimal(),
		m.err,
		continueStyle.Render("Continue ->"),
	) + "\n"
}*/
