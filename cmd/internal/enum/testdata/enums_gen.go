/**
* Generated by go-doudou v2.5.9.
* Don't edit!
 */
package testdata

import "encoding/json"

func (k *KeyboardLayout) StringSetter(value string) {
	switch value {
	case "UNKNOWN":
		*k = UNKNOWN
	case "QWERTZ":
		*k = QWERTZ
	case "AZERTY":
		*k = AZERTY
	case "QWERTY":
		*k = QWERTY
	default:
		*k = UNKNOWN
	}
}

func (k *KeyboardLayout) StringGetter() string {
	switch *k {
	case UNKNOWN:
		return "UNKNOWN"
	case QWERTZ:
		return "QWERTZ"
	case AZERTY:
		return "AZERTY"
	case QWERTY:
		return "QWERTY"
	default:
		return "UNKNOWN"
	}
}

func (k *KeyboardLayout) UnmarshalJSON(bytes []byte) error {
	var _k string
	err := json.Unmarshal(bytes, &_k)
	if err != nil {
		return err
	}
	k.StringSetter(_k)
	return nil
}

func (k *KeyboardLayout) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.StringGetter())
}
