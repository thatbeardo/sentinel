import React from "react";
import ReactMarkdown from "react-markdown";

import markdown from "../utils/docs";

const Documentation = () => (
  <div className="my-4">
    <ReactMarkdown source={markdown} />
  </div>
);

export default Documentation;
