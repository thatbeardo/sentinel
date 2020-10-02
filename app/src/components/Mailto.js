import React from "react";

const Mailto = ({ email, subject, body, wrapper, children }) => {
    return (
      <p class="hero lead mt-5 pt-5">
        {wrapper}
        <a href={`mailto:${email}?subject=${encodeURIComponent(subject) || ''}&body=${encodeURIComponent(body) || ''}`}>
            {children}</a>
      </p>
    );
  };
export default Mailto;
