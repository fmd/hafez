var gulp = require('gulp');
var sass = require('gulp-sass')
var gzip = require('gulp-gzip')

gulp.task('sass', function () {
    gulp.src('./public/scss/*.scss')
        .pipe(sass())
        .pipe(gulp.dest('./public/css'));
});

gulp.task('compress', function () {
    gulp.src('./public/css/*.css')
        .pipe(gzip())
        .pipe(gulp.dest('./public/zipped/css'));
});
