package script_vm

import (
	"blockchain_demo/pkg/sign"
	"blockchain_demo/pkg/transaction"
	"blockchain_demo/pkg/utils"
	"blockchain_demo/pkg/utils/queue"
	"blockchain_demo/pkg/utils/stack"
	"errors"
)

type operation struct{
	code OPCode
	data []byte
}
type vm struct {
	stack  stack.Stack[[]byte]
	queue  queue.Queue[operation]
	signer sign.Signer
}

type OPCode byte

const (
	OP_FALSE = 0x00 // OP_0 is the opcode for 0
	OP_TRUE  = 0x01 // OP_1 is the opcode for 1
)
const (
	OP_0            = 0x00 // OP_0 is the opcode for 0
	OP_PUSHDATA     = 0x01 // OP_1 is the opcode for 1
	OP_PUSHDATA_4B  = 0x4B // OP_1 is the opcode for 1
	OP_PUSHDATA1    = 0x4C // OP_1 is the opcode for 1
	OP_PUSHDATA2    = 0x4D // OP_1 is the opcode for 1
	OP_PUSHDATA4    = 0x4E // OP_1 is the opcode for 1
	OP_1NEGATE      = 0x4F // OP_1NEGATE is the opcode for -1

	OP_1  = 0x51 // OP_1 is the opcode for 1
	OP_2  = 0x52 // OP_2 is the opcode for 2
	OP_3  = 0x53 // OP_3 is the opcode for 3
	OP_4  = 0x54 // OP_4 is the opcode for 4
	OP_5  = 0x55 // OP_5 is the opcode for 5
	OP_6  = 0x56 // OP_6 is the opcode for 6
	OP_7  = 0x57 // OP_7 is the opcode for 7
	OP_8  = 0x58 // OP_8 is the opcode for 8
	OP_9  = 0x59 // OP_9 is the opcode for 9
	OP_10 = 0x5A // OP_10 is the opcode for 10
	OP_11 = 0x5B // OP_11 is the opcode for 11
	OP_12 = 0x5C // OP_12 is the opcode for 12
	OP_13 = 0x5D // OP_13 is the opcode for 13
	OP_14 = 0x5E // OP_14 is the opcode for 14
	OP_15 = 0x5F // OP_15 is the opcode for 15
	OP_16 = 0x60 // OP_16 is the opcode for 16

	OP_NOP    = 0x61 //	Не делает ничего	Активен
	OP_IF     = 0x63 //	Выполняет следующие утверждения если верхнее значение стека не равно 0	Активен
	OP_NOTIF  = 0x64 //	Выполняет следующие утверждения если верхнее значение стека равно 0	Активен
	OP_ELSE   = 0x67 //	Выполняет утверждения если предыдущее OP_IF или OP_NOTIF было ложным	Активен
	OP_ENDIF  = 0x68 //	Завершает блок OP_IF, OP_NOTIF или OP_ELSE	Активен
	OP_VERIFY = 0x69 //	Проверяет верхнее значение стека и прерывает, если оно равно 0	Активен
	OP_RETURN = 0x6A //	Завершает выполнение и помечает транзакцию как недействительную

	OP_IFDUP        = 0x73 //	Дублирует верхний элемент стека, если он не равен 0	Активен
	OP_DEPTH        = 0x74 //	Помещает размер стека в стек	Активен
	OP_DROP         = 0x75 //	Удаляет верхний элемент стека	Активен
	OP_DUP          = 0x76 //	Дублирует верхний элемент стека	Активен
	OP_NIP          = 0x77 //	Удаляет второй элемент стека	Активен
	OP_OVER         = 0x78 //	Копирует второй элемент стека на верх	Активен
	OP_PICK         = 0x79 //	Копирует N-й элемент стека на верх (N из вершины стека)	Активен
	OP_ROLL         = 0x7A //	Перемещает N-й элемент стека на верх	Активен
	OP_ROT          = 0x7B //	Перемещает третий элемент стека на верх	Активен
	OP_SWAP         = 0x7C //	Меняет местами два верхних элемента стека	Активен
	OP_TUCK         = 0x7D //	Копирует верхний элемент стека под второй элемент	Активен
	OP_TOALTSTACK   = 0x6B //	Перемещает верхний элемент в альтернативный стек	Активен
	OP_FROMALTSTACK = 0x6C //	Перемещает верхний элемент из альтернативного стека	Активен
	OP_2DROP        = 0x6D //	Удаляет два верхних элемента стека	Активен
	OP_2DUP         = 0x6E //	Дублирует два верхних элемента стека	Активен
	OP_3DUP         = 0x6F //	Дублирует три верхних элемента стека	Активен
	OP_2OVER        = 0x70 //	Копирует два элемента с вершин стека на верх	Активен
	OP_2ROT         = 0x71 //	Пятый и шестой элементы перемещаются на верх стека	Активен
	OP_2SWAP        = 0x72 //

	OP_BOOLAND            = 0x9A //	Логическое И двух верхних элементов стека	Активен
	OP_BOOLOR             = 0x9B //	Логическое ИЛИ двух верхних элементов стека	Активен
	OP_NUMEQUAL           = 0x9C //	Возвращает 1, если два верхних элемента равны, иначе 0	Активен
	OP_NUMEQUALVERIFY     = 0x9D //	OP_NUMEQUAL + OP_VERIFY	Активен
	OP_NUMNOTEQUAL        = 0x9E //	Возвращает 1, если два верхних элемента не равны, иначе 0	Активен
	OP_LESSTHAN           = 0x9F //	Возвращает 1, если второй сверху < верхнего, иначе 0	Активен
	OP_GREATERTHAN        = 0xA0 //	Возвращает 1, если второй сверху > верхнего, иначе 0	Активен
	OP_LESSTHANOREQUAL    = 0xA1 //	Возвращает 1, если второй сверху ≤ верхнего, иначе 0	Активен
	OP_GREATERTHANOREQUAL = 0xA2 //	Возвращает 1, если второй сверху ≥ верхнего, иначе 0	Активен
	OP_MIN                = 0xA3 //	Возвращает меньшее из двух верхних элементов	Активен
	OP_MAX                = 0xA4 //	Возвращает большее из двух верхних элементов	Активен
	OP_EQUAL              = 0x87 //	Возвращает 1, если два верхних элемента равны, иначе 0	Активен
	OP_EQUALVERIFY        = 0x88 //	OP_EQUAL + OP_VERIFY	Акти

	OP_SHA256    = 0xA8 // 	Вычисляет SHA-256 хеш верхнего элемента	Активен
	OP_HASH160   = 0xA9 // 	Вычисляет RIPEMD-160(SHA-256()) верхнего элемента	Активен
	OP_HASH256   = 0xAA // 	Вычисляет SHA-256(SHA-256()) верхнего элемента	Активен
	OP_RIPEMD160 = 0xA6 // 	Вычисляет RIPEMD-160 хеш верхнего элемента	Активен
	OP_SHA1      = 0xA7 // 	Вычисляет SHA-1 хеш верхнего элемента

	OP_CHECKSIG            = 0xAC //	Проверяет подпись и публичный ключ	Активен
	OP_CHECKSIGVERIFY      = 0xAD //	OP_CHECKSIG + OP_VERIFY	Активен
	OP_CHECKMULTISIG       = 0xAE //	Проверяет несколько подписей и публичных ключей	Активен
	OP_CHECKMULTISIGVERIFY = 0xAF //	OP_CHECKMULTISIG + OP_VERIFY
)

var ActiveCodes = [...]OPCode{
	OP_FALSE,
	OP_TRUE,
	OP_0,
	OP_1NEGATE,
	OP_1,
	OP_2, OP_3, OP_4, OP_5, OP_6, OP_7, OP_8, OP_9, OP_10, OP_11, OP_12, OP_13, OP_14, OP_15, OP_16,
	OP_PUSHDATA, OP_PUSHDATA_4B, OP_PUSHDATA1, OP_PUSHDATA2, OP_PUSHDATA4,
	OP_NOP,
	OP_IFDUP,
	OP_DROP,
	OP_DUP,
	OP_EQUAL,
	OP_EQUALVERIFY,
	OP_VERIFY,
	OP_SHA256,
	OP_HASH160,
	OP_HASH256,
	OP_CHECKSIG,
	OP_CHECKSIGVERIFY,
}
func IsActive(op OPCode) bool {
	for _, activeOp := range ActiveCodes {
		if op == activeOp {
			return true
		}
	}
	return false
}

func New(signer sign.Signer) *vm {
	return &vm{
		stack:  stack.New[[]byte](),
		signer: signer,
	}
}

func compare(a, b []byte) int {
	if len(b) > len(a) {
		return 1
	} else if len(b) < len(a) {
		return -1
	}
	for i := 0; i < len(b); i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return 0
}

func (v *vm) equal() (bool, error) {
	a, err := v.stack.Pop()
	if err != nil {
		return false, err
	}
	b, err := v.stack.Pop()
	if err != nil {
		return false, err
	}
	return compare(a, b) == 0, nil
}

func (v *vm) checksig(data []byte) (bool, error) {
	pubKey, err := v.stack.Pop()
	if err != nil {
		return false, err
	}
	signature, err := v.stack.Pop()
	if err != nil {
		return false, err
	}

	return v.signer.Verify(data, signature, pubKey)
}

func (v *vm) Precompile(script []byte) error {
	pointer := 0
	for pointer < len(script) {
		inc := 1
		dataLength := 0
		switch {
		case script[pointer] == byte(OP_0):
			v.queue.Enqueue(operation{code: OP_PUSHDATA, data: script[pointer : pointer+1]})
		case script[pointer] == OP_1NEGATE:
			v.queue.Enqueue(operation{code: OP_PUSHDATA, data: script[pointer : pointer+1]})
		case script[pointer] >= OP_1 && script[pointer] <= OP_16:
			v.queue.Enqueue(operation{code: OP_PUSHDATA, data: []byte{script[pointer] - OP_1 + 1}})
		case script[pointer] >= OP_PUSHDATA && script[pointer] <= OP_PUSHDATA_4B:
			dataLength = int(script[pointer])
			inc = dataLength + 1
			v.queue.Enqueue(operation{code: OP_PUSHDATA, data: script[pointer+1 : pointer+1+dataLength]})
		case script[pointer] == OP_PUSHDATA1:
			dataLength = int(script[pointer+1])
			inc = dataLength + 2
			v.queue.Enqueue(operation{code: OP_PUSHDATA, data: script[pointer+2 : pointer+2+dataLength]})
		case script[pointer] == OP_PUSHDATA2:
			dataLength = int(script[pointer+1]) | int(script[pointer+2])<<8
			inc = dataLength + 3
			v.queue.Enqueue(operation{code: OP_PUSHDATA, data: script[pointer+3 : pointer+3+dataLength]})
		case script[pointer] == OP_PUSHDATA4:
			dataLength = int(script[pointer+1]) | int(script[pointer+2])<<8 | int(script[pointer+3])<<16 | int(script[pointer+4])<<24
			inc = dataLength + 5
			v.queue.Enqueue(operation{code: OP_PUSHDATA, data: script[pointer+5 : pointer+5+dataLength]})
		case IsActive(OPCode(script[pointer])):
			v.queue.Enqueue(operation{code: OPCode(script[pointer]), data: nil})
		default:
			return errors.New("unknown opcode")
		}

		pointer += inc
	}
	return nil
}

func (v *vm) Run(script []byte, tx transaction.Transaction) error {
	err := v.Precompile(script)
	if err != nil {
		return err
	}
	err = v.Execute(tx)
	if err != nil {	
		return err
	}
	return nil
}

func (v *vm) Execute(tx transaction.Transaction) error {
	for op := range v.queue.Iterator() {
		switch {
		case op.code == OP_PUSHDATA:
			v.stack.Push(op.data)
		case op.code == OP_NOP:
		case op.code == OP_IFDUP:
			top, err := v.stack.Pick()
			if err != nil {
				return err
			}
			if compare(top, make([]byte, len(top))) != 0 {
				v.stack.Push(top)
			}
		case op.code == OP_DROP:
			_, err := v.stack.Pop()
			if err != nil {
				return err
			}
		case op.code == OP_DUP:
			top, err := v.stack.Pick()
			if err != nil {
				return err
			}
			v.stack.Push(top)
		case op.code == OP_EQUAL:
			eq, err := v.equal()
			if err != nil {
				return err
			}
			if eq {
				v.stack.Push([]byte{OP_TRUE})
			} else {
				v.stack.Push([]byte{OP_FALSE})
			}
		case op.code == OP_EQUALVERIFY:
			eq, err := v.equal()
			if err != nil {
				return err
			}
			if !eq {
				return errors.New("equal verify failed")
			}
		case op.code == OP_VERIFY:
			top, err := v.stack.Pop()
			if err != nil {
				return err
			}
			if compare(top, make([]byte, len(top))) == 0 {
				return errors.New("verify failed")
			}
		case op.code == OP_SHA256:
			top, err := v.stack.Pop()
			if err != nil {
				return err
			}
			hash, err := utils.GetHash(top)
			if err != nil {
				return err
			}
			v.stack.Push(hash)
		case op.code == OP_HASH160:
			top, err := v.stack.Pop()
			if err != nil {
				return err
			}
			hash, err := utils.GetHash(top)
			if err != nil {
				return err
			}
			hash160, err := utils.GetHash160(nil, hash)
			if err != nil {
				return err
			}
			v.stack.Push(hash160)
		case op.code == OP_HASH256:
			top, err := v.stack.Pop()
			if err != nil {
				return err
			}
			hash, err := utils.GetHash(top)
			if err != nil {
				return err
			}
			hash256, err := utils.GetHash(hash)
			if err != nil {
				return err
			}
			v.stack.Push(hash256)
		case op.code == OP_CHECKSIG:
			txid := tx.GetTxId()
			ok, err := v.checksig(txid[:])
			if err != nil {
				return err
			}
			if ok {
				v.stack.Push([]byte{OP_TRUE})
			} else {
				v.stack.Push([]byte{OP_FALSE})
			}
		case op.code == OP_CHECKSIGVERIFY:
			txid := tx.GetTxId()
			ok, err := v.checksig(txid[:])
			if err != nil {
				return err
			}
			if !ok {
				return errors.New("checksig verify failed")
			}
		}
	}
	return nil
}
