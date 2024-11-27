package ascii

var Art = map[string]string{
	"linux": `            ${c1}.-/+oossssoo+/-.
        ${c1}":+ssssssssssssssssss+:"
      ${c1}-+ssssssssssssssssssyyssss+-
    ${c1}.ossssssssssssssssss${c7}dMMMNy${c1}sssso.
   ${c1}/sssssssssss${c7}hdmmNNmmyNMMMMh${c1}ssssss/
  ${c1}+sssssssss${c7}h${c1}myd${c7}MMMMMMMNdddd${c1}yssssssss+
 ${c1}/ssssssssh${c7}NMM${c1}Myh${c7}hyyyyhmNMMMNh${c1}ssssssss/
${c1}.ssssssss${c7}dMMMNh${c1}ssssssssss${c7}hNMMMd${c1}ssssssss.
${c1}+ssss${c7}hhhyNMMNy${c1}ssssssssssss${c7}yNMMM${c1}ysssssss+
${c1}ossy${c7}NMMMNyMMh${c1}sssssssssssssshmmmhssssssso
${c1}ossy${c7}NMMMNyMMh${c1}sssssssssssssshmmmhssssssso
${c1}+ssss${c7}hhhyNMMNy${c1}ssssssssssss${c7}yNMMM${c1}ysssssss+
${c1}.ssssssss${c7}dMMMNh${c1}ssssssssss${c7}hNMMM${c1}dssssssss.
 ${c1}/ssssssss${c7}hNMM${c1}Myh${c7}hyyyyhdNMMMNh${c1}ssssssss/
  ${c1}+sssssssss${c7}d${c1}myd${c7}MMMMMMMMddddy${c1}ssssssss+
   ${c1}/sssssssssssh${c7}dmNNNNmyNMMMMh${c1}ssssss/
    ${c1}.ossssssssssssssssss${c7}dMMMN$${c1}ysssso.
      ${c1}-+sssssssssssssssssyyyssss+-
        ${c1}":+ssssssssssssssssss+:"
            ${c1}.-/+oossssoo+/-.

`,
	"windows": `${c1}        ,.=:!!t3Z3z.,
       :tt:::tt333EE3
${c1}       Et:::ztt33EEEL${c2} @Ee.,      ..,
${c1}      ;tt:::tt333EE7${c2} ;EEEEEEttttt33#
${c1}     :Et:::zt333EEQ.${c2} $EEEEEttttt33QL
${c1}     it::::tt333EEF${c2} @EEEEEEttttt33F
${c1}    ;3=*^""""*4EEV${c2} :EEEEEEttttt33@.
${c3}    ,.=::::!t=., ${c1}"${c2} @EEEEEEtttz33QF
${c3}   ;::::::::zt33)${c2}   "4EEEtttji3P*
${c3}  :t::::::::tt33.${c4}:Z3z..${c2}  ""${c4} ,..g.    
${c3}  i::::::::zt33F${c4} AEEEtttt::::ztF
${c3} ;:::::::::t33V${c4} ;EEEttttt::::t3
${c3} E::::::::zt33L${c4} @EEEtttt::::z3F
${c3}{3=*^""""*4E3)${c4} ;EEEtttt:::::tZ"
${c3}             "${c4} :EEEEtttt::::z7
		${c4}"VEzjt:;;z>*"
		`,
	"darwin": `${c2}                    c.'
${c2}                 ,xNMM.
${c2}               .OMMMMo
${c2}               lMMM"
${c2}     .;loddo:. .olloddol;.
${c2}   cKMMMMMMMMMMNWMMMMMMMMMM0:
${c4} .KMMMMMMMMMMMMMMMMMMMMMMMWd.
${c4} XMMMMMMMMMMMMMMMMMMMMMMMX.
${c1};MMMMMMMMMMMMMMMMMMMMMMMM:
${c1}:MMMMMMMMMMMMMMMMMMMMMMMM:
${c1}.MMMMMMMMMMMMMMMMMMMMMMMMX.
${c1} kMMMMMMMMMMMMMMMMMMMMMMMMWd.
 ${c5}'XMMMMMMMMMMMMMMMMMMMMMMMMMMk
 ${c5}'XMMMMMMMMMMMMMMMMMMMMMMMMK.
  ${c3}kMMMMMMMMMMMMMMMMMMMMMMd
   ${c3};KMMMMMMMWXXWMMMMMMMk.
    ${c3}"cooc*"    "*coo'"`,
}

var CodeColor = map[string]string{
	"reset":  "\033[0m",
	"red":    "\033[1;31m",
	"green":  "\033[1;32m",
	"cyan":   "\033[1;33m",
	"yellow": "\033[1;34m",
	"purple": "\033[1;35m",
	"blue":   "\033[1;36m",
	"white":  "\033[1;37m",
}

var PlaceHolder = map[string]string{
	"${c0}": CodeColor["reset"],
	"${c1}": CodeColor["red"],
	"${c2}": CodeColor["green"],
	"${c3}": CodeColor["yellow"],
	"${c4}": CodeColor["blue"],
	"${c5}": CodeColor["purple"],
	"${c6}": CodeColor["cyan"],
	"${c7}": CodeColor["white"],
}
