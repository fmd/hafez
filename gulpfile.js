var gulp = require('gulp');
var sass = require('gulp-sass')
var gzip = require('gulp-gzip')

gulp.task('sass', function () {
    gulp.src('./static/scss/*.scss')
        .pipe(sass())
        .pipe(gulp.dest('./public/css'));
});

gulp.tasj('less', function () {
    gulp.src('./static/less/*.less')
        .pipe(less())
        .pipe(gulp.dest('./public/css'));
});

gulp.task('compress', function () {
    gulp.src('./public/css/*.css')
        .pipe(gzip())
        .pipe(gulp.dest('./public/css'));
});
