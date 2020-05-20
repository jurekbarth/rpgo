const goLibrary = require('go-library');

const options = {
  destinationPath: 'bin',
  repo: 'jurekbarth/rpgo',
  version: 'v3.1.5',
  projectname: 'rpgo'
}


goLibrary(options);
