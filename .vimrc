autocmd FileType go nmap <leader>t  <Plug>(go-test)
autocmd BufNewFile,BufRead *.go setlocal noexpandtab tabstop=4 shiftwidth=4
Plug 'AndrewRadev/splitjoin.vim'
let g:NERDTreeDirArrowExpandable = '▸'
let g:NERDTreeDirArrowCollapsible = '▾'
if !executable('hclfmt')
  echo "'hclfmt' not found. instanlling"
  echo ""
  silent !GO111MODULE="off" go get -v github.com/fatih/hclfmt
endif
Plug 'fatih/hclfmt'
Plug 'b4b4r07/vim-hcl'
" set background=light

