
unexport VIEWDIR
unexport VIEWFILES
#
# This for AIX, we always run in 64bit mode
#
export OBJECT_MODE=64

.DEFAULT_GOAL := t40.V

t40.V: t40.v_in
	viewc -n t40.v_in


clean:
	rm *.V
