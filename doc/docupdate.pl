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


NEXT: foreach my $file (@ARGV)
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
			$lines = $lines.' '.$line;
		}
		elsif($line =~/^func/)
		{
		
			my $func = "";
			my $struct = "atmi";
			my $return = "";
			my $def = "";
			
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
			($def) = ($line =~ /(.*)\s\{/);
			
			print "func definition [$def]\n";
			
			
			if ($line =~/^func\s*\(.*\)\s*[A-Za-z0-9_]*\(.*\)\s*\(.*\)\s*{/)
			{
				# func (ac *ATMICtx) TpACall(svc string, tb TypedBuffer, flags int64) (int, ATMIError) {
				print "1: [$line]\n";
				($struct, $func, $return) = ($line =~/^func\s*\(.*\s\**(.*)\)\s*([A-Za-z0-9_]*)\(.*\)\s*\((.*)\)\s*{/);
			}
			elsif ($line =~/^func\s*\(.*\)\s*[A-Za-z0-9_]*\(.*\)\s*[0-9A-Za-z_\*]*\s*{/)
			{
				# func (ac *ATMICtx) TpACall(svc string, tb TypedBuffer, flags int64)  ATMIError {
				print "2: [$line]\n";
				($struct, $func, $return) = ($line =~/^func\s*\(.*\s\**(.*)\)\s*([A-Za-z0-9_]*)\(.*\)\s*([0-9A-Za-z_\*]*)\s*{/);
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
			
			print "func: [$func]\n";
			print "struct: [$struct]\n";
			print "return: [$return]\n";
			
			print "Building doc... (only global func, starting with capital letter)\n";
			
			if ($func!~/$[A-Z]/)
			{
				next NEXT;
			}
			
			
			my @fields = split /@/, $lines;
			
			my $descr = $fields[0];
			my $retdescr = "";
			my $varname = "";
			my $vardescr = "";
			
			$descr = $descr."\n";
			my $have_params = 0;
			
			for (my $i=1; $i < scalar @fields; $i++) {
				
				if ($fields[$i]=~/$param/)
				{
					($varname, $vardescr) = ($fields[$i] =~ /$param\s*([^\s]*)\s*(.*)/);
					
					if (!$have_params)
					{
						$descr = $descr. "*$varname* is $vardescr.";
						$have_params = 1;
					}
					else
					{
						$descr = $descr. "\n*$varname* is $vardescr.";	
					}
				}
				elsif ($fields[$i]=~/$return/)
				{
					($retdescr) = ($fields[$i] =~ /$return\s*(.*)\s*/);
				}
			}
			
			
			my $final_block = "";
			
			if ($retdescr eq "")
			{
				$final_block = <<"END_MESSAGE";
[cols="h,5a"]
|===
| Function
| $def
| Description
| $descr
|===
END_MESSAGE
			}
			else
			{
				$final_block = <<"END_MESSAGE";
[cols="h,5a"]
|===
| Function
| $def
| Description
| $descr
| Returns
| $retdescr
|===
END_MESSAGE
			}
			
			print "***************GENERATED DOC ****************\n";
			print "$final_block\n";
			print "*********************************************\n";
			
			# Link to the key
			$M_func{"$struct\.$func"} = $final_block;
		}
		else
		{
			$was_comment = 0;
			$lines = "";
		}
	
	}
	close $fh or die $!;

}


