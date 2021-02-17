import React from 'react';
import values from './values';

const CoreValues = () => (
  <div>
    <div className="text-center display-4">Our Core Values</div>
    {values.map((value, i) => (
      <div className="my-5" key={i}>
        <p className="lead">
          <strong>{value.title}</strong>
        </p>
        <p className="lead">{value.description}</p>
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
