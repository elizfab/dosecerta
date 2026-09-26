package entity

import "errors"

// Setting é um par chave/valor de configuração do usuário (app single-user).
// Ex.: chave "weightGoal" com o valor "68".
type Setting struct {
	Key   string `bson:"_id" json:"key"`
	Value string `bson:"value" json:"value"`
}

var ErrSettingKeyRequired = errors.New("chave é obrigatória")

func (s *Setting) Validate() error {
	if s.Key == "" {
		return ErrSettingKeyRequired
	}
	return nil
}
