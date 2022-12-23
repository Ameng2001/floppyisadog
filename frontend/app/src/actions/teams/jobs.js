import _ from 'lodash';
import 'whatwg-fetch';
import { normalize, schema } from 'normalizr';
import * as actionTypes from '../../constants/actionTypes';
import { routeToMicroservice } from '../../constants/paths';
import {
  emptyPromise,
  timestampExpired,
  checkStatus,
  parseJSON,
  checkCode,
} from '../../utility';

/*
  Exported functions:
  * getTeamJobs
*/

// schemas!
const teamJobsSchema = new schema.Entity('jobs', {}, { idAttribute: 'uuid' });
const arrayOfTeamJobs = new schema.Array(teamJobsSchema);

// team jobs

function requestTeamJobs(teamUuid) {
  return {
    type: actionTypes.REQUEST_TEAM_JOBS,
    teamUuid,
  };
}

function receiveTeamJobs(teamUuid, data) {
  return {
    type: actionTypes.RECEIVE_TEAM_JOBS,
    teamUuid,
    ...data,
  };
}

function fetchTeamJobs(companyUuid, teamUuid) {
  return (dispatch, getState) => {
    // dispatch action to start the fetch
    dispatch(requestTeamJobs(teamUuid));
    const jobPath = '/json/company/ListJobs';

    return fetch(
      routeToMicroservice('company', jobPath),
      { 
        //credentials: 'include',
        method: 'POST',
        mode:'cors',
        headers: {
          'Content-Type': 'application/json',
          //携带X-Token用于网关鉴权
          //'X-Token': getState().whoami.data.token,
          //whoami相当于鉴权，透传uuid和身份信息authz(设置到request header中，在网关侧会将head透传到tars调用的context中)
          'X-Verify-UID': getState().whoami.data.user_uuid,
          'X-Verify-Data': getState().whoami.data.authz,
        },
        body: JSON.stringify({
          req: {
            company_uuid: companyUuid,
            team_uuid: teamUuid
          }
        }),
      })
      .then(checkStatus)
      .then(parseJSON)
      .then(checkCode)
      .then((data) => {
        const normalized = normalize(_.get(data.rsp, 'jobs', []), arrayOfTeamJobs);

        return dispatch(receiveTeamJobs(teamUuid, {
          data: normalized.entities.jobs,
          order: normalized.result,
          lastUpdate: Date.now(),
        }));
      });
  };
}

function shouldFetchTeamJobs(state, teamUuid) {
  const jobsData = state.teams.jobs;
  const teamJobs = _.get(jobsData, teamUuid, {});

  // no team employees have ever been fetched
  if (_.isEmpty(jobsData)) {
    return true;

  // the needed teamUuid is empty
  } else if (_.isEmpty(teamJobs)) {
    return true;

  // teamJobs is at least partially populated with a trusted object at this point
  // the order of these is related to how the 1st fetch might play out

  // this data set is currently being fetched
  } else if (teamJobs.isFetching) {
    return false;

  // this data set is not complete
  } else if (!teamJobs.completeSet) {
    return true;

  // this data set is stale
  } else if (!teamJobs.lastUpdate ||
    (timestampExpired(teamJobs.lastUpdate, 'TEAM_JOBS'))
  ) {
    return true;
  }

  // check if invalidated
  return teamJobs.didInvalidate;
}

export function getTeamJobs(companyUuid, teamUuid) {
  return (dispatch, getState) => {
    if (shouldFetchTeamJobs(getState(), teamUuid)) {
      return dispatch(fetchTeamJobs(companyUuid, teamUuid));
    }
    return emptyPromise();
  };
}

export function setTeamJob(teamUuid, jobUuid, data) {
  return {
    type: actionTypes.SET_TEAM_JOB,
    teamUuid,
    jobUuid,
    data,
  };
}

function updatingTeamJob(teamUuid, jobUuid, data) {
  return {
    type: actionTypes.UPDATING_TEAM_JOB,
    teamUuid,
    jobUuid,
    data,
  };
}

function updatedTeamJob(teamUuid, jobUuid, data) {
  return {
    type: actionTypes.UPDATED_TEAM_JOB,
    teamUuid,
    jobUuid,
    data,
  };
}

function updatingTeamJobField(jobUuid) {
  return {
    type: actionTypes.UPDATING_TEAM_JOB_FIELD,
    jobUuid,
  };
}

function updatedTeamJobField(jobUuid) {
  return {
    type: actionTypes.UPDATED_TEAM_JOB_FIELD,
    jobUuid,
  };
}

function hideTeamJobFieldSuccess(jobUuid) {
  return {
    type: actionTypes.HIDE_TEAM_JOB_FIELD_SUCCESS,
    jobUuid,
  };
}

export function updateTeamJob(
  companyUuid,
  teamUuid,
  jobUuid,
  newData,
  callback
) {
  return (dispatch, getState) => {
    const jobs = _.get(getState().teams.jobs, teamUuid, {});
    const job = _.get(jobs.data, jobUuid, {});
    const updateData = _.extend({}, job, newData);
    dispatch(updatingTeamJob(teamUuid, jobUuid, newData));

    const jobPath =
      '/json/company/UpdateJob';

    return fetch(
      routeToMicroservice('company', jobPath),
      {
        //credentials: 'include',
        method: 'POST',
        mode:'cors',
        headers: {
          'Content-Type': 'application/json',
          //携带X-Token用于网关鉴权
          //'X-Token': getState().whoami.data.token,
          //whoami相当于鉴权，透传uuid和身份信息authz(设置到request header中，在网关侧会将head透传到tars调用的context中)
          'X-Verify-UID': getState().whoami.data.user_uuid,
          'X-Verify-Data': getState().whoami.data.authz,
        },
        body: JSON.stringify({
          req: {
            updateData,
          }
        }),
      })
      .then(checkStatus)
      .then(parseJSON)
      .then(checkCode)
      .then((data) => {
        if (callback) {
          callback.call(null, updateData, null);
        }

        dispatch(updatedTeamJob(teamUuid, jobUuid, updateData));
      });
  };
}

export function updateTeamJobField(companyUuid, teamUuid, jobUuid, newData) {
  return (dispatch) => {
    dispatch(updatingTeamJobField(jobUuid));

    return dispatch(
      updateTeamJob(
        companyUuid,
        teamUuid,
        jobUuid,
        newData,
        (response, error) => {
          if (!error) {
            dispatch(updatedTeamJobField(jobUuid));
            setTimeout(() => {
              dispatch(hideTeamJobFieldSuccess(jobUuid));
            }, 1000);
          }
        }
      )
    );
  };
}

function creatingTeamJob(teamUuid) {
  return {
    type: actionTypes.CREATING_TEAM_JOB,
    teamUuid,
  };
}

function createdTeamJob(teamUuid, jobUuid, data) {
  return {
    type: actionTypes.CREATED_TEAM_JOB,
    teamUuid,
    jobUuid,
    data,
  };
}

export function createTeamJob(companyUuid, teamUuid, jobPayload) {
  return (dispatch, getState) => {
    dispatch(creatingTeamJob());
    const jobPath = '/json/company/CreateJob';

    return fetch(
      routeToMicroservice('company', jobPath),
      {
        //credentials: 'include',
        method: 'POST',
        mode:'cors',
        headers: {
          'Content-Type': 'application/json',
          //携带X-Token用于网关鉴权
          //'X-Token': getState().whoami.data.token,
          //whoami相当于鉴权，透传uuid和身份信息authz(设置到request header中，在网关侧会将head透传到tars调用的context中)
          'X-Verify-UID': getState().whoami.data.user_uuid,
          'X-Verify-Data': getState().whoami.data.authz,
        },
        body: JSON.stringify({
          req: {
            jobPayload,
          }
        }),
      })
      .then(checkStatus)
      .then(parseJSON)
      .then(checkCode)
      .then((data) => {
        dispatch(createdTeamJob(teamUuid, data.rsp.uuid, data.rsp));

        setTimeout(() => {
          dispatch(hideTeamJobFieldSuccess(data.rsp.uuid));
        }, 1000);
      });
  };
}
