import React, { Fragment } from "react";
import { NavLink as RouterNavLink, useRouteMatch } from "react-router-dom";
import { Route, Switch } from "react-router-dom";
import { Nav, NavItem, NavLink, Container, Row, Col } from "reactstrap";
import Sentinel from "./Sentinel";

function Products() {
  let match = useRouteMatch();
  return (
    <Fragment>
      <Container className="mt-5">
        <Row>
          <Col xs="12">
            <Nav>
              <NavItem>
                <NavLink
                  tag={RouterNavLink}
                  to={`${match.path}/sentinel`}
                  activeClassName="router-link-exact-active"
                >
                  Sentinel
                </NavLink>
              </NavItem>
            </Nav>
          </Col>
        </Row>
        <Row>
          <Col xs="12">
            <Switch>
              <Route path={`${match.path}/sentinel`} component={Sentinel} />
            </Switch>
          </Col>
        </Row>
      </Container>
    </Fragment>
  );
}

export default Products;
