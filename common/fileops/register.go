package fileops

import "mattwach/rpngo/common/rpn"

func (fo *FileOps) register(r *rpn.RPN) {
	r.Register(".", fo.source, rpn.CatIO, sourceHelp)
	r.Register("append", fo.append, rpn.CatIO, appendHelp)
	r.Register("cd", fo.cd, rpn.CatIO, cdHelp)
	r.Register("format", fo.format, rpn.CatIO, formatHelp)
	r.Register("load", fo.load, rpn.CatIO, loadHelp)
	r.Register("save", fo.save, rpn.CatIO, saveHelp)
	r.Register("source", fo.source, rpn.CatIO, sourceHelp)
	r.Register("sh", fo.shell, rpn.CatIO, shellHelp)
}
