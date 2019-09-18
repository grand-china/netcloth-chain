// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/NetCloth/netcloth-chain/client/context
// ALIASGEN: github.com/NetCloth/netcloth-chain/client/flags
// ALIASGEN: github.com/NetCloth/netcloth-chain/client/keys
// ALIASGEN: github.com/NetCloth/netcloth-chain/client/lcd
// ALIASGEN: github.com/NetCloth/netcloth-chain/client/rest
// ALIASGEN: github.com/NetCloth/netcloth-chain/client/rpc
// ALIASGEN: github.com/NetCloth/netcloth-chain/client/tx
// ALIASGEN: github.com/NetCloth/netcloth-chain/client/utils
// ALIASGEN: github.com/NetCloth/netcloth-chain/client/input
package client

import (
	"github.com/NetCloth/netcloth-chain/client/context"
	"github.com/NetCloth/netcloth-chain/client/flags"
	"github.com/NetCloth/netcloth-chain/client/input"
	"github.com/NetCloth/netcloth-chain/client/keys"
	"github.com/NetCloth/netcloth-chain/client/lcd"
	"github.com/NetCloth/netcloth-chain/client/rpc"
)

const (
	DefaultGasAdjustment   = flags.DefaultGasAdjustment
	DefaultGasLimit        = flags.DefaultGasLimit
	GasFlagAuto            = flags.GasFlagAuto
	BroadcastBlock         = flags.BroadcastBlock
	BroadcastSync          = flags.BroadcastSync
	BroadcastAsync         = flags.BroadcastAsync
	FlagHome               = flags.FlagHome
	FlagUseLedger          = flags.FlagUseLedger
	FlagChainID            = flags.FlagChainID
	FlagNode               = flags.FlagNode
	FlagHeight             = flags.FlagHeight
	FlagGasAdjustment      = flags.FlagGasAdjustment
	FlagTrustNode          = flags.FlagTrustNode
	FlagFrom               = flags.FlagFrom
	FlagName               = flags.FlagName
	FlagAccountNumber      = flags.FlagAccountNumber
	FlagSequence           = flags.FlagSequence
	FlagMemo               = flags.FlagMemo
	FlagFees               = flags.FlagFees
	FlagGasPrices          = flags.FlagGasPrices
	FlagBroadcastMode      = flags.FlagBroadcastMode
	FlagDryRun             = flags.FlagDryRun
	FlagGenerateOnly       = flags.FlagGenerateOnly
	FlagIndentResponse     = flags.FlagIndentResponse
	FlagListenAddr         = flags.FlagListenAddr
	FlagMaxOpenConnections = flags.FlagMaxOpenConnections
	FlagRPCReadTimeout     = flags.FlagRPCReadTimeout
	FlagRPCWriteTimeout    = flags.FlagRPCWriteTimeout
	FlagOutputDocument     = flags.FlagOutputDocument
	FlagSkipConfirmation   = flags.FlagSkipConfirmation
	DefaultKeyPass         = keys.DefaultKeyPass
	FlagAddress            = keys.FlagAddress
	FlagPublicKey          = keys.FlagPublicKey
	FlagBechPrefix         = keys.FlagBechPrefix
	FlagDevice             = keys.FlagDevice
	OutputFormatText       = keys.OutputFormatText
	OutputFormatJSON       = keys.OutputFormatJSON
	MinPassLength          = input.MinPassLength
)

var (
	// functions aliases
	NewCLIContextWithFrom              = context.NewCLIContextWithFrom
	NewCLIContext                      = context.NewCLIContext
	GetFromFields                      = context.GetFromFields
	ErrInvalidAccount                  = context.ErrInvalidAccount
	ErrVerifyCommit                    = context.ErrVerifyCommit
	GetCommands                        = flags.GetCommands
	PostCommands                       = flags.PostCommands
	RegisterRestServerFlags            = flags.RegisterRestServerFlags
	ParseGas                           = flags.ParseGas
	NewCompletionCmd                   = flags.NewCompletionCmd
	MarshalJSON                        = keys.MarshalJSON
	UnmarshalJSON                      = keys.UnmarshalJSON
	Commands                           = keys.Commands
	NewAddNewKey                       = keys.NewAddNewKey
	NewRecoverKey                      = keys.NewRecoverKey
	NewUpdateKeyReq                    = keys.NewUpdateKeyReq
	NewDeleteKeyReq                    = keys.NewDeleteKeyReq
	GetKeyInfo                         = keys.GetKeyInfo
	GetPassphrase                      = keys.GetPassphrase
	ReadPassphraseFromStdin            = keys.ReadPassphraseFromStdin
	NewKeyBaseFromHomeFlag             = keys.NewKeyBaseFromHomeFlag
	NewKeyBaseFromDir                  = keys.NewKeyBaseFromDir
	NewInMemoryKeyBase                 = keys.NewInMemoryKeyBase
	NewRestServer                      = lcd.NewRestServer
	ServeCommand                       = lcd.ServeCommand
	BlockCommand                       = rpc.BlockCommand
	GetChainHeight                     = rpc.GetChainHeight
	BlockRequestHandlerFn              = rpc.BlockRequestHandlerFn
	LatestBlockRequestHandlerFn        = rpc.LatestBlockRequestHandlerFn
	RegisterRPCRoutes                  = rpc.RegisterRPCRoutes
	StatusCommand                      = rpc.StatusCommand
	NodeInfoRequestHandlerFn           = rpc.NodeInfoRequestHandlerFn
	NodeSyncingRequestHandlerFn        = rpc.NodeSyncingRequestHandlerFn
	ValidatorCommand                   = rpc.ValidatorCommand
	GetValidators                      = rpc.GetValidators
	ValidatorSetRequestHandlerFn       = rpc.ValidatorSetRequestHandlerFn
	LatestValidatorSetRequestHandlerFn = rpc.LatestValidatorSetRequestHandlerFn
	GetPassword                        = input.GetPassword
	GetCheckPassword                   = input.GetCheckPassword
	GetConfirmation                    = input.GetConfirmation
	GetString                          = input.GetString
	PrintPrefixed                      = input.PrintPrefixed

	// variable aliases
	LineBreak  = flags.LineBreak
	GasFlagVar = flags.GasFlagVar
)

type (
	CLIContext             = context.CLIContext
	GasSetting             = flags.GasSetting
	AddNewKey              = keys.AddNewKey
	RecoverKey             = keys.RecoverKey
	UpdateKeyReq           = keys.UpdateKeyReq
	DeleteKeyReq           = keys.DeleteKeyReq
	RestServer             = lcd.RestServer
	ValidatorOutput        = rpc.ValidatorOutput
	ResultValidatorsOutput = rpc.ResultValidatorsOutput
)
