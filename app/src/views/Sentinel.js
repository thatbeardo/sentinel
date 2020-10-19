import React, { Fragment } from "react";
import Hero from "../components/Hero";
import Content from "../components/Content";
import Documentation from "./Documentation";

const Sentinel = () => (
  <Fragment>
    <Hero />
    <hr />
    <Content />
    <hr />
    <Documentation/>
  </Fragment>
);

export default Sentinel;
