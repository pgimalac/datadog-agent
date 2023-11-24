# Unless explicitly stated otherwise all files in this repository are licensed
# under the Apache License Version 2.0.
# This product includes software developed at Datadog (https:#www.datadoghq.com/).
# Copyright 2016-present Datadog, Inc.

require './lib/ostools.rb'
require 'pathname'

name 'updater'

source path: '..'
relative_path 'src/github.com/DataDog/datadog-agent'

build do
  license :project_license

  # set GOPATH on the omnibus source dir for this software
  gopath = Pathname.new(project_dir) + '../../../..'
  etc_dir = "/etc/datadog-agent"
  gomodcache = Pathname.new("/modcache")
  env = {
    'GOPATH' => gopath.to_path,
    'PATH' => "#{gopath.to_path}/bin:#{ENV['PATH']}",
  }

  unless ENV["OMNIBUS_GOMODCACHE"].nil? || ENV["OMNIBUS_GOMODCACHE"].empty?
    gomodcache = Pathname.new(ENV["OMNIBUS_GOMODCACHE"])
    env["GOMODCACHE"] = gomodcache.to_path
  end

  # include embedded path (mostly for `pkg-config` binary)
  env = with_embedded_path(env)

  major_version_arg = "$MAJOR_VERSION"

  if linux_target?
    command "invoke updater.build --rebuild --major-version #{major_version_arg}", env: env
    mkdir "#{install_dir}/bin"
    mkdir "#{install_dir}/run/"


    # Config
    mkdir '/etc/datadog-agent'
    mkdir "/etc/init"
    mkdir "/var/log/datadog"

    move 'bin/agent/dist/datadog.yaml', '/etc/datadog-agent/datadog.yaml.example'
    move 'bin/agent/dist/conf.d', '/etc/datadog-agent/'
    copy 'bin/updater', "#{install_dir}/bin/"

    # Systemd
    if debian_target?
      erb source: "systemd.service.erb",
          dest: "/lib/systemd/system/datadog-agent.service",
          mode: 0644,
          vars: { install_dir: install_dir, etc_dir: etc_dir }
    else
      mkdir "/usr/lib/systemd/system/"
      erb source: "systemd.service.erb",
          dest: "/usr/lib/systemd/system/datadog-agent.service",
          mode: 0644,
          vars: { install_dir: install_dir, etc_dir: etc_dir }
    end

  end
  block do
  end

  # The file below is touched by software builds that don't put anything in the installation
  # directory (libgcc right now) so that the git_cache gets updated let's remove it from the
  # final package
  delete "#{install_dir}/uselessfile"
end
