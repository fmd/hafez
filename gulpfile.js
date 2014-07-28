var gulp = require('gulp');
var sass = require('gulp-sass')
var less = require('gulp-less')
var gzip = require('gulp-gzip')

var paths = {
  scripts: 'static/js/**/*.js',
  less:    'static/less/**/*.less',
  sass:    'static/scss/**/*.scss',
  images:  'static/images/**/*'
};

gulp.task('sass', function () {
    gulp.src(paths.sass)
        .pipe(sass())
        .pipe(gulp.dest('./public/css'));
});

gulp.task('less', function () {
    gulp.src(paths.less)
        .pipe(less())
        .pipe(gulp.dest('./public/css'));
});

gulp.task('images', function () {
    gulp.src(paths.images)
        .pipe(gulp.dest('./public/images'));
});

gulp.task('scripts', function () {
    gulp.src(paths.scripts)
        .pipe(gulp.dest('./public/js'));
});

gulp.task('watch', function() {
    gulp.watch(paths.scripts, ['scripts']);
    gulp.watch(paths.images,  ['images']);
    gulp.watch(paths.less,    ['less']);
    gulp.watch(paths.sass,    ['sass']);
});

gulp.task('default', ['sass', 'less', 'images', 'scripts', 'watch']);

/*
gulp.task('compress', function () {
    gulp.src('./public/css/*.css')
        .pipe(gzip())
        .pipe(gulp.dest('./public/css'));
});
*/

