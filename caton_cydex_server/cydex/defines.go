package cydex

const (
	OK                               = 0
	ErrUserNotExisted                = 201
	ErrInvalidPassword               = 202
	ErrLoginAgain                    = 203
	ErrInvalidUID                    = 204
	ErrRepeatedUsername              = 205
	ErrNotAllowed                    = 206
	ErrNotLogined                    = 207
	ErrFirstInvalidPassword          = 208
	ErrLicenseError                  = 209
	ErrMoreSuperUser                 = 210
	ErrKickUserFailed                = 211
	ErrBeyondUserRestriction         = 212
	ErrNoLicense                     = 213
	ErrServerNotInited               = 214
	ErrInvalidLicense                = 215
	ErrCreateGroupFailed             = 216
	ErrNoSuchGroup                   = 217
	ErrAccountExpired                = 218
	ErrOldLicenseVersion             = 219
	ErrPackageNotExisted             = 301
	ErrFileNotExisted                = 302
	ErrFullStorage                   = 304
	ErrFirewallNotEnabled            = 305
	ErrUserPathNotExisted            = 306
	ErrCreateTransferAgain           = 307
	ErrGetPortFailed                 = 308
	ErrFileNotExisted2               = 309
	ErrInvalidJsonViaAjax            = 310
	ErrSidFidNotMatch                = 311
	ErrBeyondFileRestriction         = 312
	ErrActivePackage                 = 313
	ErrUnkownLoginError              = 401
	ErrUnknowLogoutError             = 402
	ErrUnknowCreateSuperUserError    = 403
	ErrUnknowCreateUserError         = 404
	ErrUnknowModifyUserError         = 405
	ErrUnknowModifyPasswordError     = 406
	ErrUnknowQueryUserError          = 407
	ErrUnknowDelUserError            = 408
	ErrUnknowQueryLogError           = 410
	ErrUnknowGetStorageSpaceError    = 411
	ErrUnknowQueryPackageError       = 421
	ErrUnknowGetSendingPackagesError = 421
	ErrUnknowGetSendingPackageError  = 421
	ErrUnknowDelPackageError         = 421
	ErrUnknowDelFileError            = 425
	ErrUnknowCreatePackageError      = 426
	ErrUnknowCreateTransferError     = 427
	ErrInnerServer                   = 500
	ErrInvalidHttpMethod             = 501
	ErrInvalidParam                  = 502
)

const (
	UPLOAD   = 1
	DOWNLOAD = 2

	ENCRYPTION_TYPE_NONE   = 0
	ENCRYPTION_TYPE_AES128 = 1
	ENCRYPTION_TYPE_AES256 = 2

	PKG_FLAG_DELETE = 0
	PKG_FLAG_OTHER  = 1

	TRANSFER_STATE_IDLE  = 0
	TRANSFER_STATE_DOING = 1
	TRANSFER_STATE_DONE  = 2
	TRANSFER_STATE_PAUSE = 3

	FTYPE_DIR  = 0
	FTYPE_FILE = 1
)

const (
	USER_LEVEL_COMMON        = 0
	USER_LEVEL_UPLOAD_ONLY   = 1
	USER_LEVEL_DOWNLOAD_ONLY = 2
	USER_LEVEL_ADMIN         = 9
)
