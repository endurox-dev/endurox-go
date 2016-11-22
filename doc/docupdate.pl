#!/usr/bin/perl


use Text::Wrap;

#
# @(#) Read all files and group the functions by objects
# This should be run before generating a doc 
#

my $M_doc = 'endurox-go-book.adoc';

#
# So we will store Object.Func.
#
my %M_func = ();


#
# Support funcs
#

sub read_file {
    my ($filename) = @_;
 
    open my $in, '<:encoding(UTF-8)', $filename or die "Could not open '$filename' for reading $!";
    local $/ = undef;
    my $all = <$in>;
    close $in;
 
    return $all;
}
 
sub write_file {
    my ($filename, $content) = @_;
 
    open my $out, '>:encoding(UTF-8)', $filename or die "Could not open '$filename' for writing $!";;
    print $out $content;
    close $out;
 
    return;
}


$Text::Wrap::columns = 80;

foreach my $file (@ARGV)
{
	open my $fh, '<', $file or die $!;

	my $lines = "";
	my $was_comment = 0;
	my $was_func_start = 0;
NEXT:	while (<$fh>) {
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
			my $return = "N/A";
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
			elsif ($line =~/^func\s*\(.*\)\s*[A-Za-z0-9_]*\(.*\)\s*{/)
			{
				# func (ac *ATMICtx) TpACall(svc string, tb TypedBuffer, flags int64) {
				print "3: [$line]\n";
				($struct, $func) = ($line =~/^func\s*\(.*\s\**(.*)\)\s*([A-Za-z0-9_]*)\(.*\)\s*{/);
			}
			elsif ($line =~/^func\s*[A-Za-z0-9_]*\(.*\)\s*\(.*\)\s*{/)
			{
				# func NewATMICtx() (*ATMICtx, ATMIError) {
				($func, $return) = ($line =~/^func\s*([A-Za-z0-9_]*)\(.*\)\s*\((.*)\)\s*{/);
				print "4: [$line]\n";
			}
			elsif ($line =~/^func\s*[A-Za-z0-9_]*\(.*\)\s*[A-Za-z0-9_\*]*\s*{/)
			{
				# func NewATMICtx() ATMIError {
				
				($func, $return) = ($line =~/^func\s*([A-Za-z0-9_]*)\(.*\)\s*([0-9A-Za-z_\*]*)\s*{/);
				
				print "5: [$line]\n";
			}
			elsif ($line =~/^func\s*[A-Za-z0-9_]*\(.*\)\s*{/)
			{
				# func NewATMICtx() {
				
				($func) = ($line =~/^func\s*([A-Za-z0-9_]*)\(.*\)\s*{/);
				
				print "6: [$line]\n";
			}
			
			
			print "func: [$func]\n";
			print "struct: [$struct]\n";
			print "return: [$return]\n";
			
			if ($func!~/^[A-Z]/)
			{
				print "Func does not start with capital\n";
				next NEXT;
			}
			print "Building doc...\n";
			
			my @fields = split /@/, $lines;
			
			my $descr = $fields[0];
			my $retdescr = "";
			my $varname = "";
			my $vardescr = "";
			
			$descr =~ s/^\s+|\s+$//g;
			
			$descr = "$descr. ";
			#my $have_params = 0;
			
			for (my $i=1; $i < scalar @fields; $i++) {
				
				if ($fields[$i]=~/^param/)
				{
					($varname, $vardescr) = ($fields[$i] =~ /^param\s*([^\s]*)\s*(.*)\s*/);
					
					$varname =~ s/^\s+|\s+$//g;
					$vardescr =~ s/^\s+|\s+$//g;
					
					$descr = $descr."\n*$varname* $vardescr. ";
				}
				elsif ($fields[$i]=~/^return/)
				{
					($retdescr) = ($fields[$i] =~ /^return\s*(.*)\s*/);
					
					$retdescr =~ s/^\s+|\s+$//g;
				}
			}
			
			
			my $final_block = "";
			my $server_block = "";
			
			if ($file=~/atmisrv.go/)
			{
				$server_block="To XATMI server";
			}
			else
			{
				$server_block="XATMI client and server";
			}
			
			
			if ($retdescr eq "")
			{
				$final_block = <<"END_MESSAGE";
[cols="h,5a"]
|===
|Function
|$def
|Description
|$descr
|Applies
|$server_block
|===

END_MESSAGE
			}
			else
			{
				$final_block = <<"END_MESSAGE";
[cols="h,5a"]
|===
|Function
|$def
|Description
|$descr
|Returns
|$retdescr
|Applies
|$server_block
|===

END_MESSAGE
			}
			
			$final_block = wrap('', '', $final_block);
			
			print "***************GENERATED DOC ****************\n";
			print "$final_block\n";
			print "*********************************************\n";
			
			########################################################
			# set the order
			########################################################
			$prefix = "99";
			if ($struct=~/^atmi$/)
			{
				$prefix = "00";
			}
			elsif ($struct=~/^nstdError$/)
			{
				$prefix = "01";
			}
			elsif ($struct=~/^TypedJSON$/)
			{
				$prefix = "06";
			}
			elsif ($struct=~/^TypedString$/)
			{
				$prefix = "05";
			}
			elsif ($struct=~/^TypedUBF$/)
			{
				$prefix = "09";
			}
			elsif ($struct=~/^TypedCarray$/)
			{
				$prefix = "07";
			}
			elsif ($struct=~/^ATMIBuf$/)
			{
				$prefix = "03";
			}
			elsif ($struct=~/^ATMICtx$/)
			{
				$prefix = "04";
			}
			elsif ($struct=~/^atmiError$/)
			{
				$prefix = "02";
			}
			elsif ($struct=~/^ubfError$/)
			{
				$prefix = "08";
			}
		
			my $hash_key = "$prefix\.$struct\.$func";
			
			print "hash key: [$hash_key]\n";
			
			# Link to the key
			$M_func{$hash_key} = $final_block;
		}
		else
		{
			$was_comment = 0;
			$lines = "";
		}
	
	}
	close $fh or die $!;

}

### OK Seems like we got the stuff out, now need to sort and plot the doc
# https://perlmaven.com/how-to-sort-a-hash-in-perl


my $topic = "";

my $output = "";

foreach my $name (sort { $M_func{$a} <=> $M_func{$b} or $a cmp $b } keys %M_func)
{

	printf "SORTED: %-8s %s\n", $name, $M_func{$name};

	if ($name=~/^...atmi\./ && $topic ne "atmi")
	{
		$topic = "atmi";
		$output = $output."=== ATMI Package functions\n";
#
# Intro
#
my $msg = <<"END_MESSAGE";
Enduro/X package functions. ATMI Context is initiated by this package.

END_MESSAGE
		$output = $output.$msg;
	}
	elsif ($name=~/^...nstdError\./ && $topic ne "nstdError")
	{
		$topic = "nstdError";
		$output = $output."=== Enduro/X Standard Error Object / NSTDError interface\n";
#
# Intro
#
my $msg = <<"END_MESSAGE";
Enduro/X standard error object interfaced with NSTDError interface. Error is returned
by libnstd library. Which are Enduro/X base library. Currently it is used for logging.

END_MESSAGE
		$output = $output.$msg;
	}
	elsif ($name=~/^...TypedJSON\./ && $topic ne "TypedJSON")
	{
		$topic = "TypedJSON";
		$output = $output."=== JSON IPC buffer format\n";
#
# Intro
#
my $msg = <<"END_MESSAGE";
JSON buffer. Used to send JSON text between services. Basically it is string buffer,
but with special mark that it is JSON Text. This mark is special, as Enduro/X can
automatically convert JSON to UBF and vice versa. The format for JSON is one level
with UBF field names and values. Values can be arrays.

END_MESSAGE
		$output = $output.$msg;
	}
	elsif ($name=~/^...TypedString\./ && $topic ne "TypedString")
	{
		$topic = "TypedString";
		$output = $output."=== String IPC buffer format\n";
#
# Intro
#
my $msg = <<"END_MESSAGE";
String buffer. Can be used to string plain text strings between services. The string
buffer cannot contain binary zero (0x00) byte.

END_MESSAGE
		$output = $output.$msg;
	}
	elsif ($name=~/^...TypedUBF\./ && $topic ne "TypedUBF")
	{
		$topic = "TypedUBF";
		$output = $output."=== UBF Key/value IPC buffer format\n";
#
# Intro
#
my $msg = <<"END_MESSAGE";
Unified Buffer Format (UBF) is key/value buffer with compiled IDs. Each key
can contain the array of elements (occurrences).

END_MESSAGE
		$output = $output.$msg;
	}
	elsif ($name=~/^...TypedCarray\./ && $topic ne "TypedCarray")
	{
		$topic = "TypedCarray";
		$output = $output."=== Binary buffer IPC buffer format\n";
#
# Intro
#
my $msg = <<"END_MESSAGE";
Typed Carray, basically is byte array buffer.

END_MESSAGE
		$output = $output.$msg;
	}
	elsif ($name=~/^...ATMIBuf\./ && $topic ne "ATMIBuf")
	{
		$topic = "ATMIBuf";
		$output = $output."=== Abstract IPC buffer - ATMIUbf\n";
#
# Intro
#
my $msg = <<"END_MESSAGE";
ATMI buffer is base class for String, JSON, UBF (key/value with value arrays) 
and binary buffer.

END_MESSAGE
		$output = $output.$msg;
	}
	elsif ($name=~/^...ATMICtx\./ && $topic ne "ATMICtx")
	{
		$topic = "ATMICtx";
		$output = $output."=== ATMI Context\n";
#
# Intro
#
my $msg = <<"END_MESSAGE";
ATMI Context is uses as main object for accessing Enduro/X functionality. The
object is allocated by package function *atmi.NewATMICtx()*. ATMI Context API is
used for client and server API.

END_MESSAGE
		$output = $output.$msg;

	}
	elsif ($name=~/^...atmiError\./ && $topic ne "atmiError")
	{
		$topic = "atmiError";
		$output = $output."=== ATMI Error object / ATMIError interface\n";
#
# Intro
#
my $msg = <<"END_MESSAGE";
ATMI Error object, used for ATMI context functions. Error codes are described in
seperate chapter in this document.

END_MESSAGE
		$output = $output.$msg;
	}
	elsif ($name=~/^...ubfError\./ && $topic ne "ubfError")
	{
		$topic = "ubfError";
		$output = $output."=== BUF Error object/ UBFError interface\n";
#
# Intro
#
my $msg = <<"END_MESSAGE";
UBF Error object, used by UBF buffer functions.

END_MESSAGE
	}
	
	
	# Strip the leading func name
	my ($funcname) = ($name =~ /^...(.*)/);
	
	$output = $output."==== $funcname()\n";
	
	
	$output = $output.$M_func{$name}."\n";
    
}

print $output;


#
# Got to replace text between two anchors..
#
if (-e $M_doc)
{
	my $data = read_file($M_doc);
	#$data =~ s/Copyright Start-Up/Copyright Large Corporation/g;
	
	$output = "[[gen_doc-start]]\n".$output."[[gen_doc-stop]]";
	
	$data =~ s/(\[\[gen_doc-start\]\])(.*)(\[\[gen_doc-stop\]\])/$output/s;
	
	write_file($M_doc, $data);
	exit;
}
else
{
	print STDERR "$M_doc does not exists in current directory!\n";
	exit -1
}



