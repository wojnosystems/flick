package source

import (
	"github.com/wojnosystems/go-optional"
	"github.com/wojnosystems/okey-dokey/bad"
	"github.com/wojnosystems/okey-dokey/ok_int"
	"github.com/wojnosystems/okey-dokey/ok_range"
	"github.com/wojnosystems/okey-dokey/ok_string"
	"strconv"
)

type dbConfigMock struct {
	Host     optional.String
	User     optional.String
	Password optional.String
}

type dbConfigMockV struct {
	Host     ok_string.On
	User     ok_string.On
	Password ok_string.On
}

var dbConfigValidation = dbConfigMockV{
	Host: ok_string.On{
		Ensure: []ok_string.Definer{
			&ok_string.IsRequired{},
			&ok_string.LengthBetween{
				Between: ok_range.IntBetween(1, 64),
			},
		},
	},
	User: ok_string.On{
		Ensure: []ok_string.Definer{
			&ok_string.IsRequired{},
			&ok_string.LengthBetween{
				Between: ok_range.IntBetween(1, 64),
			},
		},
	},
	Password: ok_string.On{
		Ensure: []ok_string.Definer{
			&ok_string.LengthBetween{
				Between: ok_range.IntBetween(1, 64),
			},
		},
	},
}

func (v *dbConfigMockV) Validate(on *dbConfigMock, receiver bad.MemberEmitter) {
	ok_string.Validate(on.Host, &v.Host, receiver.Into("host"))
	ok_string.Validate(on.User, &v.User, receiver.Into("user"))
	ok_string.Validate(on.Password, &v.Password, receiver.Into("password"))
}

type appConfigMock struct {
	Name        optional.String
	ThreadCount optional.Int
	Databases   []dbConfigMock
}

type appConfigMockV struct {
	Name        ok_string.On
	ThreadCount ok_int.On
}

var appConfigValidation = appConfigMockV{
	Name: ok_string.On{
		Ensure: []ok_string.Definer{
			&ok_string.IsRequired{},
			&ok_string.LengthBetween{
				Between: ok_range.IntBetween(1, 25),
			},
		},
	},
	ThreadCount: ok_int.On{
		Ensure: []ok_int.Definer{
			&ok_int.IsRequired{},
			&ok_int.GreaterThanOrEqual{
				Value: 1,
			},
			&ok_int.LessThanOrEqual{
				Value: 16,
			},
		},
	},
}

func (v *appConfigMockV) Validate(on *appConfigMock, receiver bad.MemberEmitter) {
	ok_string.Validate(on.Name, &v.Name, receiver.Into("name"))
	ok_int.Validate(on.ThreadCount, &v.ThreadCount, receiver.Into("threadCount"))
	for i, database := range on.Databases {
		dbConfigValidation.Validate(&database, receiver.Into("["+strconv.FormatInt(int64(i), 10)+"]"))
	}
}

func (m dbConfigMock) IsEqual(o *dbConfigMock) bool {
	if o == nil {
		return false
	}
	return m.Host.IsEqual(o.Host) && m.User.IsEqual(o.User) && m.Password.IsEqual(o.Password)
}

func (m appConfigMock) IsEqual(o *appConfigMock) bool {
	if o == nil {
		return false
	}
	if !m.Name.IsEqual(o.Name) || !m.ThreadCount.IsEqual(o.ThreadCount) {
		return false
	}
	if len(m.Databases) != len(o.Databases) {
		return false
	}
	for i, database := range m.Databases {
		database.IsEqual(&o.Databases[i])
	}
	return true
}
