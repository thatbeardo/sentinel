import React from 'react';
import ReactMarkdown from 'react-markdown';

import './styles/styles.scss';
import gfm from 'remark-gfm';
import { resources } from "./markdown/resources";
import { contexts } from "./markdown/contexts";
import { tableOfContents } from './markdown/tableOfContents';
import { permission } from './markdown/permissions';
import { definitions } from './markdown/definitions';
import { prerequisites } from './markdown/prerequisites';
import { useCase } from './markdown/useCase';
import { creatingResources } from './markdown/creatingResources';
import { creatingContexts } from './markdown/creatingContexts';
import { creatingPermissions } from './markdown/creatingPermissions';
import { authorization } from './markdown/authorization';

const Documentation = () => (
  <div className="my-4">

    <ReactMarkdown plugins={[gfm]} children={tableOfContents} />

    <div id="resources">
      <ReactMarkdown plugins={[gfm]} children={resources} />
    </div>

    <div id="contexts">
      <ReactMarkdown plugins={[gfm]} children={contexts} />
    </div>

    <div id="permissions">
      <ReactMarkdown plugins={[gfm]} children={permission} />
    </div>

    <div id="definitions">
      <ReactMarkdown plugins={[gfm]} children={definitions} />
    </div>
  
    <div id="prerequisites">
      <ReactMarkdown plugins={[gfm]} children={prerequisites} />
    </div>

    <div id="use-case">
      <ReactMarkdown plugins={[gfm]} children={useCase} />
    </div>

    <div id="creating-resources">
      <ReactMarkdown plugins={[gfm]} children={creatingResources} />
    </div>

    <div id="creating-contexts">
      <ReactMarkdown plugins={[gfm]} children={creatingContexts} />
    </div>

    <div id="creating-permissions">
      <ReactMarkdown plugins={[gfm]} children={creatingPermissions} />
    </div>

    <div id="authorization">
      <ReactMarkdown plugins={[gfm]} children={authorization} />
    </div>
  </div>
);

export default Documentation;
