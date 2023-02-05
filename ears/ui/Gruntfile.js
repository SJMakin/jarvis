module.exports = function (grunt) {
  'use strict';

  grunt.initConfig({
    jshint: {
      sdk: {
        src: ['*.js', 'lib/*.js', 'example/*.js']
      },
      options: {
        jshintrc: '.jshintrc'
      }
    },

    copy: {
      'test-setup': {
        expand: true,
        cwd: 'test/',
        src: ['**'],
        dest: 'build/test/'
      },
      'browserify-dist-setup': {
        expand: true,
        src: ['lib/**/*.js'],
        dest: 'build/dist/',
        options: {}
      }
    },

    browserify: {
      test: {
        files: {
          'build/test/browser/lex-audio.js': 'build/dist/lib/lex-audio.js',
          'build/test/browser/specs.js': ['build/test/*-tests.js']
        }
      },
      dist: {
        files: {
          'dist/lex-audio.js': 'build/dist/lib/lex-audio.js'
        }
      }
    },

    uglify: {
      test: {
        files: {
          'build/test/browser/lex-audio.min.js': 'build/test/browser/lex-audio.js',
        }
      },
      dist: {
        files: {
          'dist/lex-audio.min.js': 'dist/lex-audio.js'
        }
      }
    },

    mocha: {
      'client-side': {
        src: ['build/test/browser/**/*.html'],
        dest: 'build/reports/client-side-tests.xml',
        options: {
          logErrors: true,
          reporter: 'XUnit',
          run: true,
          timeout: 10000,
          log: true,
          force: true
        }
      }
    },

    watch: {
      scripts: {
        files: ['lib/*.js', 'example/*.*', 'test/*tests.js', 'Gruntfile.js'],
        tasks: ['release'],
        options: {
          spawn: false,
        }
      }
    },

    clean: {
      build: {
        options: {
          force: true
        },
        src: ['./build/**']
      }
    }
  });

  grunt.loadNpmTasks('grunt-contrib-jshint');
  grunt.loadNpmTasks('grunt-contrib-copy');
  grunt.loadNpmTasks('grunt-contrib-uglify-es');
  grunt.loadNpmTasks('grunt-contrib-clean');
  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-mocha-test');
  grunt.loadNpmTasks('grunt-browserify');
  grunt.loadNpmTasks('grunt-mocha');

  grunt.registerTask('dist', ['copy:browserify-dist-setup', 'browserify:dist', 'uglify:dist']);
  grunt.registerTask('release', ['clean', 'dist']);

};