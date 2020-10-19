import React from "react";
import coreValues from "../utils/coreValues";
const CoreValues = () => (
  <div>
    <div className="text-center display-4">Our Core Values</div>

    {coreValues.map((value, i) => (
      <div className="my-5">
        <p className="lead">
          <strong>{value.title}</strong>
        </p>
        <p>{value.description}</p>
      </div>
    ))}
    <blockquote className="blockquote">
      <p className="mb-0">A house divided against itself cannot stand</p>
      <div className="blockquote-footer">
        <cite title="Source Title">Abraham Lincoln</cite>
      </div>
    </blockquote>
  </div>
);

export default CoreValues;
