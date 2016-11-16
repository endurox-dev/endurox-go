#!/usr/bin/perl

#
# @(#) Read all files and group the functions by objects
#

# https://perlmaven.com/how-to-sort-a-hash-in-perl
#foreach my $name (sort { $planets{$a} <=> $planets{$b} or $a cmp $b } keys %planets) {
#    printf "%-8s %s\n", $name, $planets{$name};
#}

#
# So we will store Object.Func.
#
%M_func = ();


foreach my $file (@ARGV)
{
	open my $fh, '<', $file or die $!;

	my $lines = "";
	my $was_comment = 0;
	my $was_func_start = 0;
	while (<$fh>) {
		# line contents's automatically stored in the $_ variable
		#print $_;
		chomp;
		
		my $line = $_;
		
		#print $line ;
		if ($line =~/^\/\//)
		{
			# Strip off the comment
			$line = substr $line, 2;
			$was_comment = 1;
			$lines = $lines.$line;
		}
		elsif($line =~/^func/)
		{
		
			if ($line=~/^func.*[^\{]\s*$/)
			{
				# We need to read some more lines here and join them
				# util we get the scope open symbol
				while (<$fh>)
				{
					chomp;
					if ($line =~/^func.*[^\{]\s*$/)
					{
						last;
					}
					$line = $line.$_;
				}
			}
			# Ok This is our func, get the structure or it is global/atmi.
			
			# func (ac *ATMICtx) TpACall(svc string, tb TypedBuffer, flags int64) ATMIError {
			print "got [$line]\n";
			if ($line =~/^func\s*\(.*\)\s*[A-Za-z0-9_]*\(.*\)\s*\(.*\)\s*{/)
			{
				# func (ac *ATMICtx) TpACall(svc string, tb TypedBuffer, flags int64) (int, ATMIError) {
				print "1: [$line]\n";
			}
			elsif ($line =~/^func\s*\(.*\)\s*[A-Za-z0-9_]*\(.*\)\s*[0-9A-Za-z_\*]*\s*{/)
			{
				# func (ac *ATMICtx) TpACall(svc string, tb TypedBuffer, flags int64)  ATMIError {
				print "2: [$line]\n";
			}
			elsif ($line =~/^func\s*[A-Za-z0-9_]*\(.*\)\s*\([0-9A-Za-z_]*\)\s*{/)
			{
				# func NewATMICtx() (*ATMICtx, ATMIError) {
				print "3: [$line]\n";
			}
			elsif ($line =~/^func\s*[A-Za-z0-9_]*\(.*\)\s*[A-Za-z0-9_\*]*s*{/)
			{
				# func NewATMICtx() ATMIError {
				print "4: [$line]\n";
			}
		}
		else
		{
			$was_comment = 0;
			$lines = "";
		}
		
		
	}
	close $fh or die $!;

}


