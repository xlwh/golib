/* module_state2.go - for collecting state info of a module  */
/*
modification history
--------------------
2014/4/24, by Zhang Miao, modify from package module_state
2014/7/9, by Li Bingyi, add SetNum feature for number states
2015/6/15, by Li Bingyi, move FormatOutput from waf-server to golang-lib
2017/12/20, by yuxiaofei, add Delete func for State
*/
/*
DESCRIPTION
This is a update version of module_state

Usage:
    import "www.baidu.com/golang-lib/module_state2"

    var state module_state2.State

    state.Init()

    state.Inc("counter", 1)
    state.Set("state", "OK")
    state.SetNum("cap", 100)
    state.SetFloat("cap", 100.1)

    stateData := state.Get()
*/
package module_state2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

import (
	"www.baidu.com/golang-lib/web_params"
)

/* state, one-level for SCounters */
type StateData struct {
	SCounters     Counters          // for count up
	States        map[string]string // for store states
	NumStates     Counters          // for store num states
	FloatStates   FloatCounters     // for store float states
	NoahKeyPrefix string            // for noah key
	ProgramName   string            // for program name
}

// state with mutex protect
type State struct {
	lock sync.Mutex
	data StateData
}

//
func NewStateData() *StateData {
	sd := new(StateData)
	sd.SCounters = NewCounters()
	sd.States = make(map[string]string)
	sd.NumStates = NewCounters()
	sd.FloatStates = NewFloatCounters()

	return sd
}

// make a copy for StateData
func (sd *StateData) copy() *StateData {
	copy := new(StateData)

	copy.SCounters = sd.SCounters.copy()

	copy.States = make(map[string]string)
	for key, value := range sd.States {
		copy.States[key] = value
	}

	copy.NumStates = NewCounters()
	for numKey, numValue := range sd.NumStates {
		copy.NumStates[numKey] = numValue
	}

	copy.FloatStates = NewFloatCounters()
	for floatKey, floatValue := range sd.FloatStates {
		copy.FloatStates[floatKey] = floatValue
	}

	copy.NoahKeyPrefix = sd.NoahKeyPrefix
	copy.ProgramName = sd.ProgramName

	return copy
}

func (sd *StateData) noahKeyGen(key string, withProgramName bool) string {
	return NoahKeyGen(key, sd.NoahKeyPrefix, sd.ProgramName, withProgramName)
}

// output noah string (lines of key:value) for StateData
func (sd *StateData) NoahString() []byte {
	return sd.noahString(false)
}

// output noah string (lines of key:value) for StateData, with program name
func (sd *StateData) NoahStringWithProgramName() []byte {
	return sd.noahString(true)
}

// escape " => \" workaround for argus collector plugin
func escapeQuote(value string) string {
	return strings.Replace(value, "\"", "\\\"", -1)
}

// output noah string (lines of key:value) for StateData
func (sd *StateData) noahString(withProgramName bool) []byte {
	var buf bytes.Buffer

	// print SCounters
	for key, value := range sd.SCounters {
		key = sd.noahKeyGen(key, withProgramName)
		str := fmt.Sprintf("%s:%d\n", key, value)
		buf.WriteString(str)
	}

	// print States
	for key, value := range sd.States {
		key = sd.noahKeyGen(key, withProgramName)
		value = escapeQuote(value)
		str := fmt.Sprintf("%s:\"%s\"\n", key, value)
		buf.WriteString(str)
	}

	// print NumStates
	for key, value := range sd.NumStates {
		key = sd.noahKeyGen(key, withProgramName)
		str := fmt.Sprintf("%s:%d\n", key, value)
		buf.WriteString(str)
	}

	// print floatStates
	for key, value := range sd.FloatStates {
		key = sd.noahKeyGen(key, withProgramName)
		str := fmt.Sprintf("%s:%f\n", key, value)
		buf.WriteString(str)
	}

	return buf.Bytes()
}

// format output according format value in params
func (sd *StateData) FormatOutput(params map[string][]string) ([]byte, error) {
	format, err := web_params.ParamsValueGet(params, "format")
	if err != nil {
		format = "json"
	}

	switch format {
	case "json":
		return json.Marshal(sd)
	case "hier_json":
		return GetSdHierJson(sd)
	case "noah":
		return sd.NoahString(), nil
	case "noah_with_program_name":
		return sd.NoahStringWithProgramName(), nil
	default:
		return nil, fmt.Errorf("format not support: %s", format)
	}
}

/* Initialize the state */
func (s *State) Init() {
	s.data.SCounters = NewCounters()
	s.data.States = make(map[string]string)
	s.data.NumStates = NewCounters()
	s.data.FloatStates = NewFloatCounters()
}

/* set noah key prefix */
func (s *State) SetNoahKeyPrefix(prefix string) {
	s.data.NoahKeyPrefix = prefix
}

/* set program name	*/
func (s *State) SetProgramName(programName string) {
	s.data.ProgramName = programName
}

/* Increase value to key */
func (s *State) Inc(key string, value int) {
	// support s is nil
	if s == nil {
		return
	}

	s.lock.Lock()
	s.data.SCounters.inc(key, value)
	s.lock.Unlock()
}

/* Decrease value to key */
func (s *State) Dec(key string, value int) {
	// support s is nil
	if s == nil {
		return
	}

	s.lock.Lock()
	s.data.SCounters.dec(key, value)
	s.lock.Unlock()
}

/* Init counters for given keys to zero */
func (s *State) CountersInit(keys []string) {
	s.lock.Lock()
	s.data.SCounters.init(keys)
	s.lock.Unlock()
}

/* set state to key */
func (s *State) Set(key string, value string) {
	// support s is nil
	if s == nil {
		return
	}

	s.lock.Lock()
	s.data.States[key] = value
	s.lock.Unlock()
}

/* delete state key */
func (s *State) Delete(key string) {
	// support s is nil
	if s == nil {
		return
	}

	s.lock.Lock()
	delete(s.data.States, key)
	s.lock.Unlock()
}

/* set num state to key */
func (s *State) SetNum(key string, value int64) {
	// support s is nil
	if s == nil {
		return
	}

	s.lock.Lock()
	s.data.NumStates[key] = value
	s.lock.Unlock()
}

/* set float state to key */
func (s *State) SetFloat(key string, value float64) {
	// support s is nil
	if s == nil {
		return
	}

	s.lock.Lock()
	s.data.FloatStates[key] = value
	s.lock.Unlock()
}

/* Get counter value of given key    */
func (s *State) GetCounter(key string) int64 {
	s.lock.Lock()
	value, ok := s.data.SCounters[key]
	s.lock.Unlock()

	if !ok {
		value = 0
	}

	return value
}

/* Get all counters */
func (s *State) GetCounters() Counters {
	s.lock.Lock()
	counters := s.data.SCounters.copy()
	s.lock.Unlock()

	return counters
}

/* Get state value of given key    */
func (s *State) GetState(key string) string {
	s.lock.Lock()
	value, ok := s.data.States[key]
	s.lock.Unlock()

	if !ok {
		value = ""
	}

	return value
}

/* Get num state value of given key    */
func (s *State) GetNumState(key string) int64 {
	s.lock.Lock()
	value, ok := s.data.NumStates[key]
	s.lock.Unlock()

	if !ok {
		value = 0
	}

	return value
}

/* Get float state value of given key    */
func (s *State) GetFloatState(key string) float64 {
	s.lock.Lock()
	value, ok := s.data.FloatStates[key]
	s.lock.Unlock()

	if !ok {
		value = float64(0.0)
	}

	return value
}

/* Get all states    */
func (s *State) GetAll() *StateData {
	s.lock.Lock()
	copy := s.data.copy()
	s.lock.Unlock()
	return copy
}

/* Get noah prefix key */
func (s *State) GetNoahKeyPrefix() string {
	return s.data.NoahKeyPrefix
}
