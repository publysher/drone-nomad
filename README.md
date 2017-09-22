# drone-nomad

Drone plugin for deploying Nomad jobs

## Example

    docker run --rm \
      -e DRONE_COMMIT_SHA=d8dbe4d94f15fe89232e0402c6e8a0ddf21af3ab \
      -e PLUGIN_JOB=jobfile.nomad \
      -e PLUGIN_USE_TEMPLATE=true \
      -e PLUGIN_NOMAD_ADDR=http://localhost:4647 \
      -v $(pwd):$(pwd) \
      -w $(pwd) \
      --privileged \
      publysher/drone-nomad