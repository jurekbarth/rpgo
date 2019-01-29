const goLibrary = require('go-library');

const options = {
  destinationPath: 'bin',
  repo: 'jurekbarth/rpgo',
  version: 'v2.0.3',
  projectname: 'rpgo'
}


goLibrary(options);
