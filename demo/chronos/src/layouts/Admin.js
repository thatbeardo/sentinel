/*!

=========================================================
* Black Dashboard React v1.1.0
=========================================================

* Product Page: https://www.creative-tim.com/product/black-dashboard-react
* Copyright 2020 Creative Tim (https://www.creative-tim.com)
* Licensed under MIT (https://github.com/creativetimofficial/black-dashboard-react/blob/master/LICENSE.md)

* Coded by Creative Tim

=========================================================

* The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

*/
import React, { useReducer } from "react";
import { Route, Switch, Redirect } from "react-router-dom";
// javascript plugin used to create scrollbars on windows
import AdminNavbar from "components/AdminNavbar";
import routes from "routes.js";

export const AppContext = React.createContext();

const initialState = {
  context: "",
};

function reducer(state, action) {
  switch (action.type) {
    case "UPDATE_CONTEXT":
      return {
        context: action.data,
      };
    default:
      return initialState;
  }
}

const getRoutes = (routes) => {
  return routes.map((prop, key) => {
    if (prop.layout === "/admin") {
      return (
        <Route
          path={prop.layout + prop.path}
          component={prop.component}
          key={key}
        />
      );
    } else {
      return null;
    }
  });
};

const Admin = (props) => {
  const [state, dispatch] = useReducer(reducer, initialState);
  return (
    <>
      <div className="wrapper">
        <div className="main-panel">
          <AppContext.Provider value={{ state, dispatch }}>
            <AdminNavbar {...props} />
            <Switch>
              {getRoutes(routes)}
              <Redirect from="*" to="/admin/dashboard" />
            </Switch>
          </AppContext.Provider>
        </div>
      </div>
    </>
  );
};

export default Admin;
