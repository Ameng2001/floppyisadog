// do not modify these - these constants are as observed in the universe
export const MILISECONDS_TO_SECONDS = 1000.0;

export const WEEK_LENGTH = 7;
export const DAYS_OF_WEEK = {
  sunday: 0,
  monday: 1,
  tuesday: 2,
  wednesday: 3,
  thursday: 4,
  friday: 5,
  saturday: 6,
};

export const UNASSIGNED_SHIFTS = 'unassigned-shifts';
export const DEFAULT_TEAM_JOB_COLOR = '673AB7';
export const NEW_JOB_UUID = 'NEW_JOB_UUID';
export const DEFAULT_NEW_JOB = {
  uuid: NEW_JOB_UUID,
  name: '',
  color: DEFAULT_TEAM_JOB_COLOR,
  isVisible: false,
};
export const COLOR_PICKER_COLORS = [
  '#F44336',
  '#E91E63',
  '#673AB7',
  '#3F51B5',
  '#2196F3',
  '#00BCD4',
  '#009688',
  '#4CAF50',
  '#FFC107',
  '#FF9800',
];

//tars return code
export const TARS_RET_OK = 100;
export const TARS_RET_Canceled = 1;
export const TARS_RET_Unknown = 2;
export const TARS_RET_InvalidArgument = 3;
export const TARS_RET_DeadlineExceeded = 4;
export const TARS_RET_NotFound = 5;
export const TARS_RET_AlreadyExists = 6;
export const TARS_RET_PermissionDenied = 7;
export const TARS_RET_ResourceExhausted = 8;
export const TARS_RET_FailedPrecondition = 9;
export const TARS_RET_Aborted = 10;
export const TARS_RET_OutOfRange = 11;
export const TARS_RET_Unimplemented = 12;
export const TARS_RET_Internal = 13;
export const TARS_RET_Unavailable = 14;
export const TARS_RET_DataLoss = 15;
export const TARS_RET_Unauthenticated = 16;

