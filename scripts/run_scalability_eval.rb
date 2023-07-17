require 'open3'
require 'date'
require 'time'

# -------------------------------- Usage --------------------------------------
#
# This script can be used to run scalability evaluations of Opera using Norma.
#
# To use the script, have Noma set up in your local repository, configure the
# parameters in the following section within this script, and run the script
# in the root directory of Norma using
#
#   /path/to/norma:> ruby ./scripts/run_scalability_eval.rb
#
# The script will run different configurations consecutively and produce a
# summary report. Paths to the report locations are printed to the console at
# the end of the evaluation run.

# ----------------------------- Configuration ---------------------------------

# The scenario to be used for the evaluation.
SCENARIO = "./scenarios/eval/scalability.yml"

# The DB implementations Opera should be evaluated with.
DB_IMPLs = [
  #"geth",
  "go-file",
]

# The number of validators to be evaluated.
NUM_VALIDATORs = [1,2,4,6,8]

# ---------------------------------- Action -----------------------------------

# Step 1 - build Norma
puts "Building ... "
build_ok = system("make -j")
if !build_ok then
    puts "Build failed, aborting."
    exit()
end
puts "OK"


# Step 2 - run Norma under various configurations
def runNorma (scenario, db, numValidators)
    puts "Running #{scenario} with #{db} and #{numValidators} validators .."
    label = "#{db}_#{numValidators}v"
    cmd = "go run ./driver/norma run --label #{label} --num-validators #{numValidators} --db-impl #{db} #{scenario}"

    puts "Running #{cmd}\n"
    
    start = Time.now    
    out = ""
    Open3.popen2e(cmd) {|stdin, stdout_and_stderr, wait_thr|
    	stdout_and_stderr.each {|line|
    	        rt = (Time.now - start).to_i
    	        rt_str = "%2d:%02d:%02d" % [rt/3600,(rt%3600)/60,rt%60]
    		puts "#{DateTime.now.strftime("%Y-%m-%d %H:%M:%S.%L")} | #{rt_str} | #{scenario} | #{db} | #{numValidators} | #{line}"
                $stdout.flush
    		out.concat(line)
    	}
    }

    res = ""
    out.scan(/Raw data was exported to (.*)/) { |file| res = file[0] }
    return res
end

$res = ["scenario, db, numValidators, measurements"]
def addResult (scenario, db, numValidators, file)
    $res.append("#{scenario}, #{db}, #{numValidators}, #{file}")
    $res.each{ |l| puts "#{l}\n" }
end

# Run the full set of configurations.
measurements = []
DB_IMPLs.each do |db|
    NUM_VALIDATORs.each do |numValidators|
        datafile = runNorma(SCENARIO, db, numValidators)
        addResult(SCENARIO, db, numValidators, datafile)
        measurements.append(datafile)
    end
end

# Step3: Generate a report summarizing the results.
cmd = "go run ./driver/norma diff #{measurements.join(" ")}"
puts "Running #{cmd} .."
system(cmd)

