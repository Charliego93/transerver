package biz

import (
	"context"
	"database/sql"
	"github.com/Charliego93/go-i18n/v2"
	"github.com/gookit/goutil/strutil"
	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nyaruka/phonenumbers"
	"github.com/transerver/accounts/internal/data/sqlc"
	"github.com/transerver/pkg/errors"
	"github.com/transerver/protos/acctspb"
	"github.com/transerver/utils"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"unicode"
)

const (
	minPasswordLength = 8
	maxPasswordLength = 32
)

type AccountRepo interface {
	ExistsByPhone(context.Context, sql.NullString) bool
	ExistsByEmail(context.Context, sql.NullString) bool
	Create(context.Context, db.AccountCreateParams) (db.Account, error)
	ByUname(context.Context, string) (db.Account, error)
}

type AccountUsecase struct {
	repo   AccountRepo
	region RegionRepo
}

func NewAccountUsecase(repo AccountRepo, region RegionRepo) *AccountUsecase {
	return &AccountUsecase{repo, region}
}

func (g *AccountUsecase) Register(ctx context.Context, req *acctspb.RegisterRequest, obj *RsaObj) error {
	var reg db.Region
	var phone, email sql.NullString
	if strutil.IsNotBlank(req.Phone) {
		req.Region = strings.ToUpper(req.Region)
		var err error
		reg, err = g.region.ByCode(context.Background(), req.Region)
		if err != nil {
			return errors.NewArgumentf(ctx, &i18n.Localized{
				MessageID:    "暂时不支持该地区",
				TemplateData: req.Region,
			})
		}

		number, err := phonenumbers.Parse(req.Phone, reg.Code)
		if err != nil {
			return errors.NewArgumentf(ctx, "手机号码和地区不匹配")
		}

		if !phonenumbers.IsValidNumberForRegion(number, req.Region) {
			return errors.NewArgumentf(ctx, "手机号码和地区不匹配")
		}

		phone = utils.SQLString(req.Phone)
		if g.repo.ExistsByPhone(ctx, phone) {
			return errors.NewArgumentf(ctx, "手机号已经存在")
		}
	} else if strutil.IsBlank(req.Email) {
		email = utils.SQLString(req.Email)
		if g.repo.ExistsByEmail(ctx, email) {
			return errors.NewArgumentf(ctx, "邮箱已经存在")
		}
	} else {
		return errors.NewArgumentf(ctx, "手机和邮箱不能同时为空")
	}

	password := utils.Bytes(req.Password)
	pwd, err := obj.Decrypt(ctx, password)
	if err != nil {
		return errors.NewInternal(ctx, "请求失败, 请刷新页面重试!")
	}

	pwdBuf, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return errors.NewInternal(ctx, "注册失败, 请尝试修改密码")
	}

	pwdLevel, err := g.passwordLevel(ctx, pwd)
	if err != nil {
		return err
	}

	account, err := g.repo.Create(ctx, db.AccountCreateParams{
		UserID:   nanoid.Must(),
		Username: req.Uname,
		Region:   req.Region,
		Area:     reg.Area,
		Phone:    phone,
		Email:    email,
		Password: pwdBuf,
		PwdLevel: int16(pwdLevel),
		Platform: "p",
	})
	_ = account
	return err
}

func (g *AccountUsecase) Login(ctx context.Context, req *acctspb.LoginRequest, obj *RsaObj) error {
	account, err := g.repo.ByUname(ctx, req.Uname)
	if err != nil {
		return err
	}
	if err != nil {
		return errors.NewArgumentf(ctx, "账户不存在")
	}

	// check req.Code
	if strutil.IsBlank(req.Code) {
		return errors.NewArgumentf(ctx, "验证码错误")
	}

	password, err := obj.Decrypt(ctx, utils.Bytes(req.Password))
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(account.Password, password)
	if err != nil {
		return errors.NewArgumentf(ctx, "密码不正确")
	}
	return nil
}

// passwordLevel 返回密码等级
// 存在特殊字符等级+1，特殊字符>5等级+1
// 存在数字等级+1，数字>5等级+1
// 存在大写字符等级+1，大写字符>5等级+1
// 存在小写字符等级+1
// 存在空白字符等级+1
func (g *AccountUsecase) passwordLevel(ctx context.Context, pwd []byte) (level uint8, err error) {
	password := []rune(utils.String(pwd))
	pwdLength := len(password)
	if pwdLength < minPasswordLength {
		return 0, errors.NewArgumentf(ctx, "密码强度过低, 不得少于?字符", minPasswordLength)
	} else if pwdLength > maxPasswordLength {
		return 0, errors.NewArgumentf(ctx, "密码过长，最长不超过?字符", maxPasswordLength)
	}

	var sc, nc, uc, lc, ec int // specialCount, numberCount, upperCount, lowerCount, spaceCount
	for _, r := range password {
		if unicode.IsControl(r) {
			return 0, errors.NewArgumentf(ctx, "密码包含非法字符")
		}

		if unicode.IsUpper(r) {
			uc++
		}
		if unicode.IsLower(r) {
			lc++
		}
		if unicode.IsSymbol(r) || unicode.IsPunct(r) || unicode.IsLetter(r) {
			sc++
		}
		if unicode.IsSpace(r) {
			ec++
		}
		if unicode.IsNumber(r) {
			nc++
		}
	}

	if nc == 0 {
		return 0, errors.NewArgumentf(ctx, "密码必须包含数字")
	} else {
		level++
	}
	if uc == 0 {
		return 0, errors.NewArgumentf(ctx, "密码必须包含大写字母")
	} else {
		level++
	}
	if lc == 0 {
		return 0, errors.NewArgumentf(ctx, "密码必须包含小写字母")
	} else {
		level++
	}

	if sc > 0 {
		level++
		if sc > 5 {
			level++
		}
	}
	if nc > 5 {
		level++
	}
	if uc > 5 {
		level++
	}
	if ec > 0 {
		level++
	}
	return level, nil
}
