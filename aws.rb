require 'aws-sdk'

Aws.config.update({
  region: ENV['REGION'],
  credentials: Aws::Credentials.new(
    ENV['ACCESS_KEY'],
    ENV['SECRET_KEY']
  )
})

s3 = Aws::S3::Resource.new
obj = s3.bucket('go-s3-learning').object('ruby-sdk-testfile')
p obj.upload_file('testfile')
