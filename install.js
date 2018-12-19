const goLibrary = require('go-library');

const options = {
  destinationPath: 'bin',
  repo: 'jurekbarth/rpgo',
  version: 'v2.0.1',
  projectname: 'rpgo'
}


goLibrary(options);
