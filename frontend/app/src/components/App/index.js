import _ from 'lodash';
import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import NavigationSide from 'components/SideNavigation';
import Intercom from 'components/Intercom';
import * as paths from 'constants/paths';
import * as actions from 'actions';

import Drawer, { DrawerContent, DrawerAppContent } from '@rmwc/drawer';

require('./app.scss');


class App extends Component {
  componentDidMount() {
    const { dispatch, companyUuid } = this.props;

    // query whoami endpoint if needed
    dispatch(actions.getWhoAmI());

    // get user data too
    dispatch(actions.getUser());

    // get company info because we are now at the company level
    dispatch(actions.getCompany(companyUuid));

    // get team data because it's needed for side nav paths
    dispatch(actions.getTeams(companyUuid));

    // get intercom settings
    dispatch(actions.fetchIntercomSettings());
  }

  render() {
    const { children, companyUuid, intercomSettings } = this.props;

    return (
      <div className='drawer-container'>
        <NavigationSide companyUuid={companyUuid} />

        <DrawerAppContent className='drawer-app-content'>
          {children}
        </DrawerAppContent>

        {!_.isEmpty(intercomSettings)
        &&
          <Intercom
            {...intercomSettings}
            appID={intercomSettings.app_id}
          />}
      </div>
    );
  }
}

App.propTypes = {
  dispatch: PropTypes.func.isRequired,
  children: PropTypes.element,
  companyUuid: PropTypes.string.isRequired,
  intercomSettings: PropTypes.object.isRequired,
};

function mapStateToProps(state, ownProps) {
  return {
    companyId: ownProps.routeParams.companyId,
    intercomSettings: state.whoami.intercomSettings,
  };
}

export default connect(mapStateToProps)(App);
