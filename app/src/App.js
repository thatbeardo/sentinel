import React from "react";
import { Router, Route, Switch, Redirect } from "react-router-dom";
import { Container } from "reactstrap";

import PrivateRoute from "./components/PrivateRoute";
import Loading from "./components/Loading";
import NavBar from "./components/NavBar";
import Footer from "./components/Footer";
import AboutUs from "./views/AboutUs";
import Products from "./views/Products";
import PrivacyPolicy from "./views/PrivacyPolicy";
import TermsOfUse from "./views/TermsOfUse";
import ContactUs from "./views/ContactUs";
import Documentation from "./views/Documentation";
import Profile from "./views/Profile";
import { useAuth0 } from "./react-auth0-spa";
import history from "./utils/history";
import ReactGA from 'react-ga';

// styles
import "./App.css";

// fontawesome
import initFontAwesome from "./utils/initFontAwesome";
initFontAwesome();

const trackingId = "G-CE9JDJWW4F"; 
const App = () => {

  useEffect(() => {
    ReactGA.initialize(trackingId);
    ReactGA.pageView('/')
  }, [])

  const { loading } = useAuth0();

  if (loading) {
    return <Loading />;
  }

  return (
    <Router history={history}>
      <div id="app" className="d-flex flex-column h-100">
        <NavBar />
        <Container className="flex-grow-1">
          <Switch>
            <Route path="/aboutus" exact component={AboutUs} />
            <Route path="/contactus" exact component={ContactUs} />
            <Route path="/documentation" exact component={Documentation} />
            <PrivateRoute path="/profile" component={Profile} />
            <Route path="/products" component={Products} />
            <Route path="/terms-of-use" component={TermsOfUse} />
            <Route path="/privacy-policy" component={PrivacyPolicy} />
            <Route path="/terms" component={TermsOfUse} />
            <Route path="/privacy" component={PrivacyPolicy} />
            <Route path="/">
              <Redirect to="products/sentinel" />
            </Route>
          </Switch>
        </Container>
        <Footer />
      </div>
    </Router>
  );
};

export default App;
