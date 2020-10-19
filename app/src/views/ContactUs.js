import React, { Fragment } from "react";

import Mailto from "../components/Mailto";

const ContactUs = () => (
  <Fragment>
    <Mailto
      email="hmavani7@gmail.com"
      subject="Hello, Harshil!"
      body="I would like to learn more about Sentinel!"
      wrapper="Interested to learn more? Drop us "
      children="an email"
    ></Mailto>
  </Fragment>
);

export default ContactUs;
