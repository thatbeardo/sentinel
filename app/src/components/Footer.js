import React from "react";

const Footer = () => (
  <footer className="bg-light footer">
    <div className="container">      
      <div className="grid">
        <div className="row">
          <a href="/privacy">Privacy Policy </a>
          &nbsp;&nbsp;
          <a href="/terms">Terms of Use</a>
        </div>

        <div className="row">
          &copy; {new Date().getFullYear()} Copyright BitHippie. All Rights
          Reserved.
        </div>
      </div>
    </div>
  </footer>
);

export default Footer;
