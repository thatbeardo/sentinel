import React from "react";

import logo from "../assets/sentinel.png";

const Hero = () => (
  <div className="text-center hero">
    <img className="sentinel-logo" src={logo} alt="React logo" />

    <p className="lead mb-5">
      A domain agnostic, application layer authorization solution.
    </p>
  </div>
);

export default Hero;
